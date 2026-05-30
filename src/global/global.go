package global

import (
	"gin-admin/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Config *config.Config
	Log    *zap.Logger
)
