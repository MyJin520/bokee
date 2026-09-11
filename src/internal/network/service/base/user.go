package base

import (
	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/pkg/cryptox/hash"
	"bokee/pkg/jwtx"
	"bokee/pkg/redisx"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserService struct{}

func (s *UserService) Create(req request.UserCreateReq) error {
	var count int64
	global.DB.Model(&basic.User{}).Where("phone = ? OR email = ?", req.Phone, req.Email).Count(&count)
	if count > 0 {
		global.Log.Warn("注册失败：手机号或邮箱已被注册",
			zap.String("phone", req.Phone), zap.String("email", req.Email))
		return fmt.Errorf("注册失败>手机号或邮箱已被注册")
	}

	hashedPwd, err := hash.GeneratePassword(req.Password)
	if err != nil {
		global.Log.Error("用户密码加密失败", zap.Error(err))
		return fmt.Errorf("注册失败，密码加密失败请联系管理员")
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
		return fmt.Errorf("注册失败，请稍后重试")
	}
	return nil
}

func (s *UserService) Login(req request.UserLoginReq) (*response.JwtResp, error) {
	var user basic.User
	err := global.DB.Select("id", "password", "name", "status", "avatar", "email", "phone").Preload("Roles").
		Where("name = ?", req.Username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("登录失败：用户名不存在", zap.String("username", req.Username))
			return nil, fmt.Errorf("用户名或密码错误")
		}
		global.Log.Error("登录数据库查询失败", zap.Error(err), zap.String("username", req.Username))
		return nil, fmt.Errorf("系统繁忙，请稍后重试")
	}

	if user.Status != "normal" {
		global.Log.Warn("账号已被禁用", zap.String("username", req.Username))
		return nil, fmt.Errorf("账号已被禁用，请联系管理员")
	}

	if err := hash.CompareHashAndPassword(user.Password, req.Password); err != nil {
		global.Log.Warn("密码校验失败", zap.String("username", req.Username))
		return nil, fmt.Errorf("用户名或密码错误")
	}

	roleCodes := make([]uint, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleCodes = append(roleCodes, role.RoleCode)
	}

	token, err := jwtx.GenerateToken(user.ID, user.Name, roleCodes)
	if err != nil {
		global.Log.Error("生成Token失败", zap.Error(err), zap.Uint("userID", user.ID))
		return nil, fmt.Errorf("登录失败，请稍后重试")
	}

	jwtResponse := &response.JwtResp{
		Token:  token,
		Name:   user.Name,
		Avatar: user.Avatar,
		Email:  user.Email,
		Phone:  user.Phone,
	}
	return jwtResponse, nil
}

func (s *UserService) Update(req request.UserUpdateReq, uid uint) error {
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
		global.Log.Warn("更新用户失败：没有需要更新的字段", zap.Uint("userID", uid))
		return fmt.Errorf("没有需要更新的字段")
	}

	result := global.DB.Model(&basic.User{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		global.Log.Error("更新用户失败", zap.Error(result.Error), zap.Uint("userID", uid))
		return fmt.Errorf("更新用户失败，请稍后重试")
	}
	if result.RowsAffected == 0 {
		global.Log.Warn("更新用户失败：用户不存在或未做任何更改", zap.Uint("userID", uid))
		return fmt.Errorf("用户不存在或未做任何更改")
	}
	return nil
}

// Logout 用户登出：将当前 token 加入 Redis 黑名单
func (s *UserService) Logout(c *gin.Context) error {
	// 从请求头提取 token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		global.Log.Warn("登出失败：未提供认证令牌")
		return fmt.Errorf("未提供认证令牌")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		global.Log.Warn("登出失败：认证令牌格式错误")
		return fmt.Errorf("认证令牌格式错误")
	}
	tokenString := parts[1]

	// 解析 token 获取过期时间
	claims, err := jwtx.ParseToken(tokenString)
	if err != nil {
		global.Log.Warn("登出失败：无效的认证令牌", zap.Error(err))
		return fmt.Errorf("无效的认证令牌")
	}

	// 计算 token 剩余有效期，作为黑名单 TTL
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		return nil // token 已过期，无需处理
	}

	// 写入 Redis 黑名单
	if err := redisx.BlacklistToken(c, tokenString, remaining); err != nil {
		global.Log.Warn("登出写入 Redis 黑名单失败", zap.Error(err))
		return fmt.Errorf("登出失败，请稍后重试")
	}

	global.Log.Info("用户登出成功",
		zap.Uint("userID", claims.UserID),
		zap.String("username", claims.Username),
		zap.Duration("blacklist_ttl", remaining),
	)
	return nil
}

func (s *UserService) ForgetPassword(req request.ForgetPasswordReq, cruId uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 查询当前操作用户并判断是否为管理员
		var currentUser basic.User
		err := tx.Where("id = ?", cruId).Preload("Roles").First(&currentUser).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				global.Log.Warn("重置密码失败：用户不存在", zap.Uint("userID", cruId))
				return fmt.Errorf("用户不存在")
			}
			global.Log.Error("查询用户失败", zap.Error(err), zap.Uint("userID", cruId))
			return fmt.Errorf("查询用户失败，请稍后重试")
		}

		isAdmin := false
		for _, role := range currentUser.Roles {
			if role.RoleCode == global.SuperRoleCode {
				isAdmin = true
				break
			}
		}

		// 确定目标用户 ID
		var targetUserID uint
		if isAdmin {
			if req.SpecifyUserID != 0 {
				targetUserID = req.SpecifyUserID
			} else {
				targetUserID = cruId
			}
		} else {
			if req.SpecifyUserID != 0 && req.SpecifyUserID != cruId {
				global.Log.Warn("重置密码失败：非管理员不能修改其他用户的密码", zap.Uint("userID", cruId))
				return fmt.Errorf("非管理员不能修改其他用户的密码")
			}
			targetUserID = cruId
		}

		// 非管理员校验旧密码（currentUser 已包含 password，无需二次查询）
		if !isAdmin {
			if err := hash.CompareHashAndPassword(currentUser.Password, req.OldPassword); err != nil {
				global.Log.Warn("重置密码失败：旧密码错误", zap.Uint("userID", cruId))
				return fmt.Errorf("旧密码错误，请输入正确的旧密码")
			}
		}

		// 管理员跨用户操作时，确认目标用户存在（轻量 Count 代替全量查询）
		if isAdmin && targetUserID != cruId {
			var exists int64
			if err := tx.Model(&basic.User{}).Where("id = ?", targetUserID).Count(&exists).Error; err != nil {
				global.Log.Error("查询目标用户失败", zap.Error(err), zap.Uint("targetUserID", targetUserID))
				return fmt.Errorf("查询目标用户失败，请稍后重试")
			}
			if exists == 0 {
				global.Log.Warn("重置密码失败：目标用户不存在", zap.Uint("targetUserID", targetUserID))
				return fmt.Errorf("目标用户不存在")
			}
		}

		// 加密新密码
		newPassword, err := hash.GeneratePassword(req.NewPassword)
		if err != nil {
			global.Log.Error("密码加密失败", zap.Error(err))
			return fmt.Errorf("密码加密失败，请稍后重试")
		}

		// 更新密码
		if err := tx.Model(&basic.User{}).Where("id = ?", targetUserID).Update("password", newPassword).Error; err != nil {
			global.Log.Error("更新密码失败", zap.Error(err), zap.Uint("targetUserID", targetUserID))
			return fmt.Errorf("修改密码失败，请稍后重试")
		}

		global.Log.Info("密码修改成功",
			zap.Uint("operatorID", cruId),
			zap.Uint("targetUserID", targetUserID),
			zap.Bool("isAdmin", isAdmin))
		return nil
	})
}

// BindRoles 为用户绑定角色（追加式：在已有角色基础上追加指定角色，不影响已绑定的角色）
func (s *UserService) BindRoles(req request.UserRoleBindReq) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 1. 查询用户是否存在
		var user basic.User
		err := tx.Where("id = ?", req.UserID).First(&user).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				global.Log.Warn("绑定角色失败：用户不存在", zap.Uint("userID", req.UserID))
				return fmt.Errorf("用户不存在")
			}
			global.Log.Error("查询用户失败", zap.Error(err), zap.Uint("userID", req.UserID))
			return fmt.Errorf("查询用户失败，请稍后重试")
		}

		// 2. 校验所有角色是否存在
		if len(req.RoleIDs) == 0 {
			global.Log.Warn("绑定角色失败：角色ID列表不能为空", zap.Uint("userID", req.UserID))
			return fmt.Errorf("角色ID列表不能为空")
		}

		var roles []basic.Role
		if err := tx.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
			global.Log.Error("查询角色失败", zap.Error(err))
			return fmt.Errorf("查询角色失败，请稍后重试")
		}

		if len(roles) != len(req.RoleIDs) {
			global.Log.Warn("绑定角色失败：部分角色不存在", zap.Uint("userID", req.UserID), zap.Any("roleIDs", req.RoleIDs))
			return fmt.Errorf("部分角色不存在，请检查角色ID")
		}

		// 3. 获取用户已绑定的角色，过滤出尚未绑定的角色进行追加
		var currentRoles []basic.Role
		if err := tx.Model(&user).Association("Roles").Find(&currentRoles); err != nil {
			global.Log.Error("查询用户当前角色失败", zap.Error(err))
			return fmt.Errorf("查询用户当前角色失败，请稍后重试")
		}

		currentRoleIDs := make(map[uint]bool)
		for _, r := range currentRoles {
			currentRoleIDs[r.ID] = true
		}

		var newRoles []basic.Role
		for _, r := range roles {
			if !currentRoleIDs[r.ID] {
				newRoles = append(newRoles, r)
			}
		}

		if len(newRoles) == 0 {
			global.Log.Warn("绑定角色失败：指定角色已绑定", zap.Uint("userID", req.UserID))
			return fmt.Errorf("指定角色已绑定，无需重复绑定")
		}

		// 4. 追加新角色（GORM many2many 自动写入 sys_user_roles 表）
		if err := tx.Model(&user).Association("Roles").Append(&newRoles); err != nil {
			global.Log.Error("绑定用户角色失败", zap.Error(err),
				zap.Uint("userID", req.UserID), zap.Any("roleIDs", req.RoleIDs))
			return fmt.Errorf("角色绑定失败，请稍后重试")
		}

		global.Log.Info("用户角色绑定成功",
			zap.Uint("userID", req.UserID),
			zap.Any("roleIDs", req.RoleIDs),
			zap.Int("newRoleCount", len(newRoles)),
		)
		return nil
	})
}

func (s *UserService) List(req request.UserListReq) ([]basic.User, int64, error) {
	query := global.DB.Model(&basic.User{})

	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.Email != "" {
		query = query.Where("email LIKE ?", "%"+req.Email+"%")
	}
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		global.Log.Error("统计用户总数失败", zap.Error(err))
		return nil, 0, fmt.Errorf("查询用户列表失败，请稍后重试")
	}

	var users []basic.User
	err = query.Offset(req.Offset()).Limit(req.PageSize).Preload("Roles").Find(&users).Error
	if err != nil {
		global.Log.Error("查询用户列表失败", zap.Error(err))
		return nil, 0, fmt.Errorf("查询用户列表失败，请稍后重试")
	}

	return users, total, nil
}

func (s *UserService) GetInfo(id uint) (basic.User, error) {
	var user basic.User
	err := global.DB.Where("id = ?", id).Preload("Roles").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("获取用户信息失败：用户不存在", zap.Uint("userID", id))
			return basic.User{}, fmt.Errorf("用户不存在")
		}
		global.Log.Error("查询用户失败", zap.Error(err), zap.Uint("userID", id))
		return basic.User{}, fmt.Errorf("查询用户失败，请稍后重试")
	}
	return user, nil
}