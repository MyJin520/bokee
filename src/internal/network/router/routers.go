package router

import (
	"gin-admin/internal/network/middleware"
	"gin-admin/internal/network/router/base"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	engine := gin.Default()

	// 公开路由组
	publicGroup := engine.Group("/pub")

	// 私有路由组
	privateGroup := engine.Group("/pri")
	privateGroup.Use(middleware.AuthMiddleware())
	privateGroup.Use(middleware.CasbinMiddleware())

	base.InitUserRouter(publicGroup, privateGroup)

	return engine
}
