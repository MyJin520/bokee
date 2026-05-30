package initialize

import (
	"gin-admin/global"
)

func Init() {
	var err error
	global.Config, err = initLoadConfig("C:\\KimJin\\work\\ideaProject\\gin-admin\\src\\config.yaml")
	if err != nil {
		panic("初始化配置文件失败: " + err.Error())
	}
	initLogger(global.Config)
	initDb(global.Config)
}
