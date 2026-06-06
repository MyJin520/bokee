package base

import (
	"gin-admin/internal/network/api/base"
	"github.com/gin-gonic/gin"
)

var roleApi = base.RoleApi{}

func InitRoleRouter(privateGroup *gin.RouterGroup) {
	rolePrivate := privateGroup.Group("/role")
	{
		rolePrivate.GET("/get_all_pri_rule", roleApi.GetAllPriRoles) // 获取所有私密路由
	}
}
