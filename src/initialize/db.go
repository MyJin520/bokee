package initialize

import (
	"gin-admin/config"
	db "gin-admin/pkg/gormx"
)

func initDb(cfg *config.Config) {
	err := db.InitDB(cfg)
	if err != nil {
		panic("数据库初始化失败: " + err.Error())
	}
}
