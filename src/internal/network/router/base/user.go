package base

import (
	"gin-admin/internal/network/api/base"
	"github.com/gin-gonic/gin"
)

var userApi = base.UserApi{}

func InitUserRouter(publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	userPublic := publicGroup.Group("/user")
	{
		userPublic.POST("/register", userApi.Register)
		userPublic.POST("/login", userApi.Login)
	}

	userPrivate := privateGroup.Group("/user")
	{
		userPrivate.POST("/edit", userApi.Edit)
		userPrivate.GET("/logout", userApi.Logout)
		userPrivate.POST("/token_parsing", userApi.ParseToken)
	}
}
