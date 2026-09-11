package base

import (
	"bokee/global"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/internal/network/service/base"
	"bokee/pkg/verifyx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type RoleApi struct{}

var roleService = &base.RoleService{}

// Create 创建角色
func (a *RoleApi) Create(c *gin.Context) {
	var req request.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("创建角色请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("创建角色参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}

	resp, err := roleService.Create(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(resp, "创建角色成功", c)
}

// Update 更新角色
func (a *RoleApi) Update(c *gin.Context) {
	var req request.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("更新角色请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("更新角色参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
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
		global.Log.Warn("角色ID格式错误", zap.String("id", idStr))
		response.FailWithRequest("无效的角色ID", c)
		return
	}

	if err := roleService.Delete(uint(id)); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("删除角色成功", c)
}

// GetInfo 获取角色详情
func (a *RoleApi) GetInfo(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		global.Log.Warn("角色ID格式错误", zap.String("id", idStr))
		response.FailWithRequest("无效的角色ID", c)
		return
	}

	resp, err := roleService.GetInfo(uint(id))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(resp, "获取角色详情成功", c)
}

// List 分页获取角色列表
func (a *RoleApi) List(c *gin.Context) {
	var req request.RoleQueryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("角色列表请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("角色列表参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
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
		global.Log.Error("角色授权请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("角色授权参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}

	if err := roleService.Auth(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("角色授权成功", c)
}

// GetAllPriRule 获取所有私有路由（保留原有功能）
func (a *RoleApi) GetAllPriRule(c *gin.Context) {
	var page request.PageReq
	if err := c.ShouldBindJSON(&page); err != nil {
		global.Log.Error("私有路由列表请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	page.Normalize()

	routes, total, err := roleService.GetAllPriRule(page)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(routes, total, page.Page, page.PageSize, "获取私有路由列表成功", c)
}
