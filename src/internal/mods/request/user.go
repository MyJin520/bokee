package request

// UserCreateReq 创建用户
type UserCreateReq struct {
	Username string `json:"username" label:"用户名" validate:"required"`
	Password string `json:"password" label:"密码" validate:"required,min=6"`
	Phone    string `json:"phone" label:"手机号"`
	Email    string `json:"email" label:"邮箱" validate:"omitempty,email"`
	Avatar   string `json:"avatar"`
}

// UserLoginReq 用户登陆
type UserLoginReq struct {
	Username string `json:"username" label:"用户名" validate:"required"`
	Password string `json:"password" label:"密码" validate:"required"`
}

// TokenParsingReq todo 测试接口后续删除
type TokenParsingReq struct {
	Token string `json:"token" label:"Token" validate:"required"`
}

// UserUpdateReq 更新用户
type UserUpdateReq struct {
	Username string `json:"username" label:"用户名"`
	Phone    string `json:"phone" label:"手机号"`
	Email    string `json:"email" label:"邮箱" validate:"omitempty,email"`
	Avatar   string `json:"avatar"`
}

// ForgetPasswordReq 忘记密码
type ForgetPasswordReq struct {
	OldPassword   string `json:"old_password" label:"旧密码" validate:"required"`
	NewPassword   string `json:"new_password" label:"新密码" validate:"required,min=6"`
	SpecifyUserID uint   `json:"specify_user_id"`
}

// UserListReq 用户列表查询请求
type UserListReq struct {
	PageReq
	Name   string `json:"name"`   // 用户名模糊搜索
	Phone  string `json:"phone"`  // 手机号模糊搜索
	Email  string `json:"email"`  // 邮箱模糊搜索
	Status string `json:"status"` // 账号状态精确匹配
}

// UserRoleBindReq 用户角色绑定请求（支持批量绑定多个角色）
type UserRoleBindReq struct {
	UserID  uint   `json:"userId" label:"用户ID" validate:"required"`          // 用户ID
	RoleIDs []uint `json:"roleIds" label:"角色ID列表" validate:"required,min=1"` // 角色ID列表（传单个角色ID也使用数组格式）
}
