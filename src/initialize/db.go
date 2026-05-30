package initialize

import (
	"gin-admin/config"
	"gin-admin/global"
	"gin-admin/internal/mods/basic"
	db "gin-admin/pkg/gormx"
)

func initDb(cfg *config.Config) {
	err := db.InitDB(cfg)
	if err != nil {
		panic("数据库初始化失败: " + err.Error())
	}
	migrateTable()
}

// 迁移表结构
func migrateTable() {
	err := global.DB.AutoMigrate(
		// 基础模型
		&basic.User{},
		&basic.Role{},
	)
	if err != nil {
		panic("数据库迁移失败: " + err.Error())
	}
}
