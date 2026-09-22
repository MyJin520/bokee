package basic

type Files struct {
	BaseModel
	Url              string `gorm:"size:125;not null;comment:文件地址" json:"url"`
	Ext              string `gorm:"size:10;not null;comment:文件后缀" json:"ext"`
	Size             int64  `gorm:"not null;comment:文件大小" json:"size"`
	OriginalFileName string `gorm:"size:125;not null;comment:原始文件名" json:"originalFileName"`
	Hash             string `gorm:"size:255;not null;uniqueIndex;comment:文件hash" json:"hash"`
}
