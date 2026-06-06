package base

import (
	"gin-admin/internal/network/api/base"
	"github.com/gin-gonic/gin"
)

var roleApi = base.RoleApi{}

func InitRoleRouter(privateGroup *gin.RouterGroup) {
	rolePrivate := privateGroup.Group("/role")
	{
		rolePrivate.POST("/create", roleApi.Create)                  // 创建角色
		rolePrivate.PUT("/update", roleApi.Update)                   // 更新角色
		rolePrivate.DELETE("/delete", roleApi.Delete)                // 删除角色
		rolePrivate.GET("/get_info", roleApi.Get)                    // 获取角色详情
		rolePrivate.GET("/list", roleApi.List)                       // 分页获取角色列表
		rolePrivate.GET("/get_all_pri_rule", roleApi.GetAllPriRoles) // 获取所有私有路由（原有功能）
	}
}
