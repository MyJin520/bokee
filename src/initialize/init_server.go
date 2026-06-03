package initialize

import (
	"gin-admin/global"
)

func Init() {
	var err error
	global.Config, err = initLoadConfig("./config.yaml")
	if err != nil {
		panic("初始化配置文件失败: " + err.Error())
	}
	initLogger(global.Config)
	initRedis(global.Config)
	initDb(global.Config)
	initCasbin()
	ginServerInit(global.Config.System)
}
