package base

import (
	"gin-admin/internal/network/api/base"
	"github.com/gin-gonic/gin"
)

var userApi = base.UserApi{}

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("/user")
	{
		userRouter.POST("register", userApi.Register)        // 注册用户
		userRouter.POST("login", userApi.Login)              // 登录用户
		userRouter.POST("token_parsing", userApi.ParseToken) // 解析Token
	}
}
