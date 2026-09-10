package router

import (
	"bokee/internal/network/middleware"
	"bokee/internal/network/router/articles"
	"bokee/internal/network/router/base"
	"bokee/internal/network/router/files"
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
	articles.InitArticleRouter(privateGroup)
	files.InitFileRouter(publicGroup)

	return engine
}
