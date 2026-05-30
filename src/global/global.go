package global

import (
	"gin-admin/config"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Config *config.Config
)
