package base

import (
	"errors"
	"gin-admin/global"
	"gin-admin/internal/mods/basic"
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/pkg/crypto/hash"
	"gin-admin/pkg/jwtx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct{}

func (s *UserService) Register(req request.UserRegisterOrEditReq) error {
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
