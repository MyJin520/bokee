package request

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	RoleName string `json:"roleName" label:"角色名称" validate:"required"`
	RoleCode uint   `json:"roleCode" label:"角色标识" validate:"required"`
	Sort     int    `json:"sort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	ID       uint   `json:"id" label:"角色ID" validate:"required"`
	RoleName string `json:"roleName" label:"角色名称"`
	Sort     int    `json:"sort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// RoleQueryReq 角色分页查询请求
type RoleQueryReq struct {
	RoleName string `json:"roleName"`
	Status   string `json:"status"`
	PageReq         // 嵌入分页参数
}

// RoleAuthItem 单条授权规则
type RoleAuthItem struct {
	Path   string `json:"path" label:"请求路径" validate:"required"`   // 请求路径
	Method string `json:"method" label:"请求方法" validate:"required"` // 请求方法（GET/POST/PUT/DELETE）
}

// RoleAuthReq 角色授权请求（支持批量）
type RoleAuthReq struct {
	RoleID uint           `json:"roleID" label:"角色ID" validate:"required"`             // 角色ID
	Rules  []RoleAuthItem `json:"rules" label:"授权规则列表" validate:"required,min=1,dive"` // 授权规则列表
}
