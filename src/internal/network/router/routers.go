package router

import (
	"gin-admin/internal/network/router/base"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	engine := gin.Default()

	// 公开路由组
	publicGroup := engine.Group("/pub")
	{
		base.InitUserRouter(publicGroup)
	}

	return engine
}
