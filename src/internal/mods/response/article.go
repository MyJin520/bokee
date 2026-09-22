package response

import "time"

// AuthorInfo 文章作者公开信息（脱敏：仅展示必要字段）
type AuthorInfo struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ArticleInfoResp 文章详情响应结构体
type ArticleInfoResp struct {
	ID        uint       `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Summary   string     `json:"summary"`
	Cover     string     `json:"cover"`
	ViewCount int        `json:"viewCount"`
	LikeCount int        `json:"likeCount"`
	IsTop     bool       `json:"isTop"`
	UserID    uint       `json:"userId"`
	Author    AuthorInfo `json:"author"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// ArticleListItemResp 文章列表项响应结构体
type ArticleListItemResp struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Cover     string    `json:"cover"`
	ViewCount int       `json:"viewCount"`
	LikeCount int       `json:"likeCount"`
	IsTop     bool      `json:"isTop"`
	UserID    uint      `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
