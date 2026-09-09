package base

import (
	"bokee/internal/network/api/base"
	"github.com/gin-gonic/gin"
)

var roleApi = base.RoleApi{}

func InitRoleRouter(privateGroup *gin.RouterGroup) {
	rolePrivate := privateGroup.Group("/role")
	{
		rolePrivate.POST("/create", roleApi.Create)                  // 创建角色
		rolePrivate.PUT("/update", roleApi.Update)                   // 更新角色
		rolePrivate.DELETE("/delete", roleApi.Delete)                // 删除角色
		rolePrivate.GET("/get_info", roleApi.GetInfo)                // 获取角色详情
		rolePrivate.POST("/list", roleApi.List)                      // 分页获取角色列表
		rolePrivate.POST("/auth", roleApi.Auth)                      // 角色授权（批量）
		rolePrivate.POST("/get_all_pri_rule", roleApi.GetAllPriRule) // 获取所有私有路由（原有功能）
	}
}
