package router

import (
	"gin-admin/internal/network/middleware"
	"gin-admin/internal/network/router/base"
	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	engine := gin.Default()

	// CORS 跨域中间件（全局）
	engine.Use(middleware.CORSMiddleware())

	// 公开路由组
	publicGroup := engine.Group("/pub")

	// 私有路由组
	privateGroup := engine.Group("/pri")
	privateGroup.Use(middleware.AuthMiddleware())
	privateGroup.Use(middleware.CasbinMiddleware())

	base.InitUserRouter(publicGroup, privateGroup)
	base.InitRoleRouter(privateGroup)

	return engine
}
