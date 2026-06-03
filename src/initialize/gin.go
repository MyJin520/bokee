package initialize

import (
	"gin-admin/config"
	"gin-admin/internal/network/router"
	"gin-admin/pkg/casbinx"
)

func ginServerInit(systemConfig config.SystemConfig) {
	engine := router.Routers()
	casbinx.InitSuperRoleCasbin(engine)
	err := engine.Run(systemConfig.Addr)
	if err != nil {
		panic("服务启动异常> " + err.Error())
	}
}
