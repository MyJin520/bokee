package response

import "time"

type JwtResp struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

// UserRoleResp 用户角色响应结构体（脱敏：仅暴露展示所需字段）
type UserRoleResp struct {
	ID       uint   `json:"id"`
	RoleName string `json:"roleName"`
	RoleCode uint   `json:"roleCode"`
}

// UserInfoResp 用户信息响应结构体（脱敏：剔除密码等敏感字段）
type UserInfoResp struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Status    string         `json:"status"`
	Avatar    string         `json:"avatar"`
	Roles     []UserRoleResp `json:"roles"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

type PriRouteResp struct {
	Path   string `json:"path"`   // 路由路径
	Method string `json:"method"` // 请求方法
}
