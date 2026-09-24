package base

import (
	"bokee/internal/network/api/base"
	"bokee/pkg/routex"
	"github.com/gin-gonic/gin"
)

var userActionApi = base.UserActionApi{}

func InitUserActionRouter(privateGroup *gin.RouterGroup) {
	userActionPrivate := routex.NewGroup("用户操作模块", privateGroup.Group("/user_action"))
	{
		userActionPrivate.POST("/create", "创建用户操作", userActionApi.Create)
	}
}
