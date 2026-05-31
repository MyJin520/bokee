package base

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"gin-admin/global"
	"gin-admin/internal/mods/basic"
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/jwtx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserService struct{}

func (s *UserService) Register(req request.UserRegisterReq) error {
	var count int64
	global.DB.Model(&basic.User{}).Where("phone = ? OR email = ?", req.Phone, req.Email).Count(&count)
	if count > 0 {
		return errors.New("注册失败>手机号或邮箱已被注册")
	}

	hashedPwd, err := hash.GeneratePassword(req.Password)
	if err != nil {
		global.Log.Error("用户密码加密失败", zap.Error(err))
		return errors.New("注册失败，密码加密失败请联系管理员")
	}

	user := &basic.User{
		Name:     req.Username,
		Password: hashedPwd,
		Phone:    req.Phone,
		Email:    req.Email,
		Avatar:   req.Avatar,
	}

	if err := global.DB.Create(user).Error; err != nil {
		global.Log.Error("注册失败", zap.Error(err))
		return errors.New("注册失败，请稍后重试")
	}
	return nil
}

func (s *UserService) Login(req request.UserLoginReq) (*response.JwtResponse, error) {
	var user basic.User
	err := global.DB.Select("id", "password", "name", "status", "avatar", "email", "phone").
		Where("name = ?", req.Username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("登录失败：用户名不存在", zap.String("username", req.Username))
			return nil, errors.New("用户名或密码错误")
		}
		global.Log.Error("登录数据库查询失败", zap.Error(err), zap.String("username", req.Username))
		return nil, errors.New("系统繁忙，请稍后重试")
	}

	if user.Status != "normal" {
		global.Log.Warn("账号已被禁用", zap.String("username", req.Username))
		return nil, errors.New("账号已被禁用，请联系管理员")
	}

	if err := hash.CompareHashAndPassword(user.Password, req.Password); err != nil {
		global.Log.Warn("密码校验失败", zap.String("username", req.Username))
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwtx.GenerateToken(user.ID, user.Name)
	if err != nil {
		global.Log.Error("生成Token失败", zap.Error(err), zap.Uint("userID", user.ID))
		return nil, errors.New("登录失败，请稍后重试")
	}

	jwtResponse := &response.JwtResponse{
		Token:  token,
		Name:   user.Name,
		Avatar: user.Avatar,
		Email:  user.Email,
		Phone:  user.Phone,
	}
	return jwtResponse, nil
}

func (s *UserService) Edit(req request.UserEditReq, uid uint) error {
	// 构建需要更新的字段映射，只更新非空字段
	updates := make(map[string]interface{})

	if req.Username != "" {
		updates["name"] = req.Username
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		// TODO: 验证邮箱验证码后再允许更新
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	if len(updates) == 0 {
		return errors.New("没有需要更新的字段")
	}

	result := global.DB.Model(&basic.User{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在或未做任何更改")
	}
	return nil
}

// blacklistKey 生成 token 黑名单 key（SHA256 哈希，兼容旧 token）
func blacklistKey(tokenString string) string {
	h := sha256.Sum256([]byte(tokenString))
	return "gin-admin:token:blacklist:" + fmt.Sprintf("%x", h)
}

// Logout 用户登出：将当前 token 加入 Redis 黑名单
func (s *UserService) Logout(c *gin.Context) error {
	// 1. 从请求头提取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return errors.New("未提供认证令牌")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return errors.New("认证令牌格式错误")
	}
	tokenString := parts[1]

	// 2. 解析 token 获取过期时间
	claims, err := jwtx.ParseToken(tokenString)
	if err != nil {
		return errors.New("无效的认证令牌")
	}

	// 3. 计算 token 剩余有效期，作为黑名单 TTL
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		return nil // token 已过期，无需处理
	}

	// 4. 写入 Redis 黑名单（使用 SHA256 哈希做 key，不依赖 jti）
	if global.Redis != nil {
		key := blacklistKey(tokenString)
		err := global.Redis.Set(c, key, "1", remaining).Err()
		if err != nil {
			global.Log.Warn("登出写入 Redis 黑名单失败", zap.Error(err))
			return errors.New("登出失败，请稍后重试")
		}
	}

	global.Log.Info("用户登出成功",
		zap.Uint("userID", claims.UserID),
		zap.String("username", claims.Username),
		zap.Duration("blacklist_ttl", remaining),
	)
	return nil
}
