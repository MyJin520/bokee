package request

type UserRegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type UserLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenParsingReq struct {
	Token string `json:"token"`
}

type UserEditReq struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type ForgetPasswordReq struct {
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password"`
	SpecifyUserID uint   `json:"specify_user_id"`
}

// UserRoleBindReq 用户角色绑定请求（支持批量绑定多个角色）
type UserRoleBindReq struct {
	UserID  uint   `json:"userId" binding:"required"`  // 用户ID
	RoleIDs []uint `json:"roleIds" binding:"required"` // 角色ID列表（传单个角色ID也使用数组格式）
}
