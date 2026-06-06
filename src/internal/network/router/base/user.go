package base

import (
	"gin-admin/internal/network/api/base"
	"github.com/gin-gonic/gin"
)

var userApi = base.UserApi{}

func InitUserRouter(publicGroup *gin.RouterGroup, privateGroup *gin.RouterGroup) {
	userPublic := publicGroup.Group("/user")
	{
		userPublic.POST("/register", userApi.Register) // 注册用户
		userPublic.POST("/login", userApi.Login)       // 用户登录
	}

	userPrivate := privateGroup.Group("/user")
	{
		userPrivate.POST("/edit", userApi.Edit)                      // 更新用户信息
		userPrivate.GET("/logout", userApi.Logout)                   // 用户登出
		userPrivate.GET("/get_info", userApi.GetInfo)                // 获取用户信息
		userPrivate.POST("/forget_password", userApi.ForgetPassword) // 忘记密码
	}
	{
		userPrivate.GET("/get_all_pri_rule", userApi.GetAllPriRoles) // 获取所有私密路由
	}
}
