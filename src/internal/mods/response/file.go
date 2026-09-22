package response

import "time"

// FileResp 文件上传响应结构体
type FileResp struct {
	ID               uint      `json:"id"`
	Url              string    `json:"url"`
	Ext              string    `json:"ext"`
	Size             int64     `json:"size"`
	OriginalName     string    `json:"originalName"`
	Hash             string    `json:"hash"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}