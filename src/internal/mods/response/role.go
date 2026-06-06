package response

import "time"

// RoleResp 角色响应结构体
type RoleResp struct {
	ID        uint      `json:"id"`
	RoleName  string    `json:"roleName"`
	RoleCode  uint      `json:"roleCode"`
	Sort      int       `json:"sort"`
	Status    string    `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
