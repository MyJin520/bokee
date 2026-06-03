package global

import (
	"gin-admin/config"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
	Log    *zap.Logger
)

const (
	SuperRoleCode = 888
)
