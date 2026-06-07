package base

import (
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/internal/network/service/base"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RoleApi struct{}

var roleService = &base.RoleService{}

// Create 创建角色
func (a *RoleApi) Create(c *gin.Context) {
	var req request.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}

	resp, err := roleService.Create(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(resp, "创建角色成功", c)
}

// Update 更新角色
func (a *RoleApi) Update(c *gin.Context) {
	var req request.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}

	if err := roleService.Update(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新角色成功", c)
}

// Delete 删除角色
func (a *RoleApi) Delete(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(http.StatusBadRequest, "无效的角色ID", c)
		return
	}

	if err := roleService.Delete(uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除角色成功", c)
}

// Get 获取角色详情
func (a *RoleApi) Get(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Fail(http.StatusBadRequest, "无效的角色ID", c)
		return
	}

	resp, err := roleService.Get(uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(resp, "获取角色详情成功", c)
}

// List 分页获取角色列表
func (a *RoleApi) List(c *gin.Context) {
	var req request.RoleQueryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	req.Normalize()

	list, total, err := roleService.List(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(list, total, req.Page, req.PageSize, "获取角色列表成功", c)
}

// Auth 角色授权
func (a *RoleApi) Auth(c *gin.Context) {
	var req request.RoleAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}

	if err := roleService.Auth(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("角色授权成功", c)
}

// GetAllPriRoles 获取所有私有路由（保留原有功能）
func (a *RoleApi) GetAllPriRoles(c *gin.Context) {
	var page request.PageReq
	if err := c.ShouldBindJSON(&page); err != nil {
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
