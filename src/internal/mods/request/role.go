package request

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	RoleName string `json:"roleName"`
	RoleCode uint   `json:"roleCode"`
	Sort     int    `json:"sort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	ID       uint   `json:"id"`
	RoleName string `json:"roleName"`
	Sort     int    `json:"sort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// RoleQueryReq 角色分页查询请求
type RoleQueryReq struct {
	RoleName string `form:"roleName"`
	Status   string `form:"status"`
	PageReq         // 嵌入分页参数
}
