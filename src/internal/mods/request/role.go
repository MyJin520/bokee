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

// RoleAuthItem 单条授权规则
type RoleAuthItem struct {
	Path   string `json:"path"`   // 请求路径
	Method string `json:"method"` // 请求方法（GET/POST/PUT/DELETE）
}

// RoleAuthReq 角色授权请求（支持批量）
type RoleAuthReq struct {
	RoleID uint           `json:"roleID" binding:"required"` // 角色ID
	Rules  []RoleAuthItem `json:"rules" binding:"required"`  // 授权规则列表
}
