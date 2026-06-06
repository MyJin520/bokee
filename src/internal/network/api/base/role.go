package base

import (
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/internal/network/service/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RoleApi struct{}

var roleService = &base.RoleService{}

// GetAllPriRoles todo 仅管理员访问能访问
func (u *RoleApi) GetAllPriRoles(c *gin.Context) {
	var page request.PageReq
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	page.Normalize()

	routes, total, err := roleService.GetAllPriRoles(page)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(routes, total, page.Page, page.PageSize, "获取私有路由列表成功", c)
}
