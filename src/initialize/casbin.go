package initialize

import (
	"gin-admin/pkg/casbinx"
)

func initCasbin() {
	if err := casbinx.InitEnforcer(); err != nil {
		panic("Casbin 权限系统初始化失败: " + err.Error())
	}
}
