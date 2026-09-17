package base

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bokee/global"
	"bokee/internal/mods/basic"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/pkg/cryptox/hash"
	"bokee/pkg/jwtx"
	"bokee/pkg/limitx"
	"bokee/pkg/redisx"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService struct{}

var (
	// loginLimiter 登录防爆破：10min 窗口内 5 次失败锁定 30min
	loginLimiter = &limitx.FailLimiter{Scope: "user:login", Window: 10 * time.Minute, Max: 5, Lock: 30 * time.Minute}
	// resetLimiter 重置密码防爆破：10min 窗口内 5 次失败锁定 30min（仅非管理员旧密码错误计数）
	resetLimiter = &limitx.FailLimiter{Scope: "user:reset", Window: 10 * time.Minute, Max: 5, Lock: 30 * time.Minute}
)

// userInfoKey 用户信息缓存 Key：bokee:user:info:<id>
func userInfoKey(id uint) string {
	return redisx.BuildKey("user", "info", fmt.Sprintf("%d", id))
}

// userRegisterKey 注册限流 Key：bokee:user:register:<phone>
func userRegisterKey(phone string) string {
	return redisx.BuildKey("user", "register", phone)
}

// deleteUserInfoCache 写操作落库成功后失效用户信息缓存（先写库、后删缓存，保证最终一致）
func deleteUserInfoCache(ctx context.Context, ids ...uint) {
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, userInfoKey(id))
	}
	if err := redisx.Delete(ctx, keys...); err != nil {
		global.Log.Error("删除用户信息缓存失败", zap.Strings("keys", keys), zap.Error(err))
	}
}

// buildUserInfoResp 将用户模型转换为脱敏响应结构（剔除密码等敏感字段）
func buildUserInfoResp(user basic.User) response.UserInfoResp {
	roles := make([]response.UserRoleResp, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, response.UserRoleResp{
			ID:       role.ID,
			RoleName: role.RoleName,
			RoleCode: role.RoleCode,
		})
	}
	return response.UserInfoResp{
		ID:        user.ID,
		Name:      user.Name,
		Phone:     user.Phone,
		Email:     user.Email,
		Status:    user.Status,
		Avatar:    user.Avatar,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// Create 用户注册：内置注册防刷（同一手机号窗口内仅允许尝试一次），Redis 异常放行不影响注册
func (s *UserService) Create(ctx context.Context, req request.UserCreateReq) error {
	// 注册防刷：SetNX 成功才继续，失败说明窗口内已尝试过（Redis 异常放行）
	if ok, err := redisx.SetNX(ctx, userRegisterKey(req.Phone), "1", 1*time.Minute); err != nil {
		global.Log.Error("注册限流检查失败", zap.Error(err), zap.String("phone", req.Phone))
	} else if !ok {
		global.Log.Warn("注册过于频繁", zap.String("phone", req.Phone))
		return fmt.Errorf("注册过于频繁，请稍后再试")
	}

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

// Login 用户登录：内置登录防爆破（失败计数 + 锁定），Redis 异常放行不影响正常登录
func (s *UserService) Login(ctx context.Context, req request.UserLoginReq) (*response.JwtResp, error) {
	// 登录防爆破：先检查是否已被锁定（Redis 异常放行）
	if remain, locked := loginLimiter.IsLocked(ctx, req.Username); locked {
		global.Log.Warn("登录失败：账号已被锁定", zap.String("username", req.Username))
		return nil, fmt.Errorf("尝试次数过多，请%d分钟后再试", int(remain.Minutes())+1)
	}

	var user basic.User
	err := global.DB.Select("id", "password", "name", "status", "avatar", "email", "phone").Preload("Roles").
		Where("name = ?", req.Username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("登录失败：用户名不存在", zap.String("username", req.Username))
			loginLimiter.RecordFail(ctx, req.Username)
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
		loginLimiter.RecordFail(ctx, req.Username)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 登录成功：清除失败计数与锁定标记
	loginLimiter.Clear(ctx, req.Username)

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

// Update 更新当前用户资料，落库成功后失效用户信息缓存
func (s *UserService) Update(ctx context.Context, req request.UserUpdateReq, uid uint) error {
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

	// 落库成功后失效缓存
	deleteUserInfoCache(ctx, uid)
	return nil
}

// Logout 用户登出：将当前 token 加入 Redis 黑名单（token 由 API 层从请求头提取后传入，不依赖 gin）
func (s *UserService) Logout(ctx context.Context, tokenString string) error {
	// 解析 token 获取过期时间
	claims, err := jwtx.ParseToken(tokenString)
	if err != nil {
		global.Log.Error("登出失败：无效的认证令牌", zap.Error(err))
		return fmt.Errorf("无效的认证令牌")
	}

	// 计算 token 剩余有效期，作为黑名单 TTL
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining <= 0 {
		return nil // token 已过期，无需处理
	}

	// 写入 Redis 黑名单
	if err := redisx.BlacklistToken(ctx, tokenString, remaining); err != nil {
		global.Log.Error("登出写入 Redis 黑名单失败", zap.Error(err))
		return fmt.Errorf("登出失败，请稍后重试")
	}

	global.Log.Info("用户登出成功",
		zap.Uint("userID", claims.UserID),
		zap.String("username", claims.Username),
		zap.Duration("blacklist_ttl", remaining),
	)
	return nil
}

// ForgetPassword 重置密码：内置防爆破（非管理员旧密码错误计数 + 锁定），Redis 异常放行不影响主流程
func (s *UserService) ForgetPassword(ctx context.Context, req request.ForgetPasswordReq, currentUserID uint) error {
	// 重置密码防爆破：先检查当前操作者是否已被锁定（Redis 异常放行）
	operatorID := fmt.Sprintf("%d", currentUserID)
	if remain, locked := resetLimiter.IsLocked(ctx, operatorID); locked {
		global.Log.Warn("重置密码失败：操作已被锁定", zap.Uint("userID", currentUserID))
		return fmt.Errorf("尝试次数过多，请%d分钟后再试", int(remain.Minutes())+1)
	}

	// 目标用户 ID（事务内确定，成功后用于失效缓存）
	var targetUserID uint
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		// 查询当前操作用户并判断是否为管理员
		var currentUser basic.User
		err := tx.Where("id = ?", currentUserID).Preload("Roles").First(&currentUser).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				global.Log.Warn("重置密码失败：用户不存在", zap.Uint("userID", currentUserID))
				return fmt.Errorf("用户不存在")
			}
			global.Log.Error("查询用户失败", zap.Error(err), zap.Uint("userID", currentUserID))
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
		if isAdmin {
			if req.SpecifyUserID != 0 {
				targetUserID = req.SpecifyUserID
			} else {
				targetUserID = currentUserID
			}
		} else {
			if req.SpecifyUserID != 0 && req.SpecifyUserID != currentUserID {
				global.Log.Warn("重置密码失败：非管理员不能修改其他用户的密码", zap.Uint("userID", currentUserID))
				return fmt.Errorf("非管理员不能修改其他用户的密码")
			}
			targetUserID = currentUserID
		}

		// 非管理员校验旧密码（currentUser 已包含 password，无需二次查询）
		if !isAdmin {
			if err := hash.CompareHashAndPassword(currentUser.Password, req.OldPassword); err != nil {
				global.Log.Warn("重置密码失败：旧密码错误", zap.Uint("userID", currentUserID))
				resetLimiter.RecordFail(ctx, operatorID)
				return fmt.Errorf("旧密码错误，请输入正确的旧密码")
			}
		}

		// 管理员跨用户操作时，确认目标用户存在（轻量 Count 代替全量查询）
		if isAdmin && targetUserID != currentUserID {
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
			zap.Uint("operatorID", currentUserID),
			zap.Uint("targetUserID", targetUserID),
			zap.Bool("isAdmin", isAdmin))
		return nil
	})
	if err != nil {
		return err
	}

	// 重置成功后清除失败计数与锁定标记，并失效目标用户信息缓存
	resetLimiter.Clear(ctx, operatorID)
	deleteUserInfoCache(ctx, targetUserID)
	return nil
}

// BindRoles 为用户绑定角色（追加式：在已有角色基础上追加指定角色，不影响已绑定的角色），成功后失效用户信息缓存
func (s *UserService) BindRoles(ctx context.Context, req request.UserRoleBindReq) error {
	err := global.DB.Transaction(func(tx *gorm.DB) error {
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
	if err != nil {
		return err
	}

	// 角色变更成功：失效该用户的信息缓存（缓存含角色信息）
	deleteUserInfoCache(ctx, req.UserID)
	return nil
}

// List 分页获取用户列表（返回脱敏结构，剔除密码等敏感字段）
func (s *UserService) List(req request.UserListReq) ([]response.UserInfoResp, int64, error) {
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

	var modelUsers []basic.User
	err = query.Offset(req.Offset()).Limit(req.PageSize).Preload("Roles").Find(&modelUsers).Error
	if err != nil {
		global.Log.Error("查询用户列表失败", zap.Error(err))
		return nil, 0, fmt.Errorf("查询用户列表失败，请稍后重试")
	}

	users := make([]response.UserInfoResp, 0, len(modelUsers))
	for _, u := range modelUsers {
		users = append(users, buildUserInfoResp(u))
	}
	return users, total, nil
}

// GetInfo 获取用户信息
func (s *UserService) GetInfo(ctx context.Context, id uint) (response.UserInfoResp, error) {
	key := userInfoKey(id)

	var cached response.UserInfoResp
	if err := redisx.GetJSON(ctx, key, &cached); err != nil {
		global.Log.Error("读取用户信息缓存失败，降级查库", zap.Error(err), zap.Uint("userID", id))
	} else if cached.ID != 0 {
		return cached, nil
	}

	var user basic.User
	if err := global.DB.Preload("Roles").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Log.Warn("获取用户信息失败：用户不存在", zap.Uint("userID", id))
			return response.UserInfoResp{}, fmt.Errorf("用户不存在")
		}
		global.Log.Error("查询用户失败", zap.Error(err), zap.Uint("userID", id))
		return response.UserInfoResp{}, fmt.Errorf("查询用户失败，请稍后重试")
	}

	resp := buildUserInfoResp(user)
	if err := redisx.SetJSON(ctx, key, resp, 30*time.Minute); err != nil {
		global.Log.Error("回填用户信息缓存失败", zap.Error(err), zap.Uint("userID", id))
	}
	return resp, nil
}
