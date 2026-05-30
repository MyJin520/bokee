package base

import (
	"gin-admin/internal/network/api/base"
	"github.com/gin-gonic/gin"
)

var userApi = base.UserApi{}

func InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("/user")
	{
		userRouter.POST("register", userApi.Register) // 注册用户
	}
}
