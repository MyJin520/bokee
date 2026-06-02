package initialize

import (
	"gin-admin/config"
	"gin-admin/internal/network/router"
)

func ginServerInit(systemConfig config.SystemConfig) {
	engine := router.Routers()
	err := engine.Run(systemConfig.Addr)
	if err != nil {
		panic("服务启动异常> " + err.Error())
	}
}
