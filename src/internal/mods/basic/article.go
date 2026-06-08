package basic

type Article struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" comment:"作者ID"`
	Title     string `gorm:"type:varchar(200);not null" comment:"文章标题"`
	Content   string `gorm:"type:longtext;not null" comment:"正文内容(Markdown或HTML)"`
	Summary   string `gorm:"type:varchar(300)" comment:"文章摘要"`
	Cover     string `gorm:"type:varchar(500)" comment:"封面图URL"`
	ViewCount int    `gorm:"default:0" comment:"阅读量"`
	LikeCount int    `gorm:"default:0" comment:"点赞数"`
	IsTop     bool   `gorm:"default:false;index" comment:"是否置顶"`
}
