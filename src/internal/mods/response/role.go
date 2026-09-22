package response

import "time"

// RoleResp 角色响应结构体
type RoleResp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      uint      `json:"code"`
	Sort      int       `json:"sort"`
	Status    string    `json:"status"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}