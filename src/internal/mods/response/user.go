package response

type JwtResp struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
}

// RoleItemResp  角色列表项
type RoleItemResp struct {
	ID       uint   `json:"id"`
	RoleName string `json:"roleName"`
	RoleCode uint   `json:"roleCode"`
	Sort     int    `json:"sort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}
