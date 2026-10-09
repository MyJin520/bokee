package base

import (
	"bokee/global"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/internal/network/middleware"
	"bokee/internal/network/service/base"
	"bokee/pkg/routex"
	"bokee/pkg/verifyx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserActionApi struct{}

var userActionService = &base.UserActionService{}

func (u *UserActionApi) Create(c *gin.Context) {
	var req request.ActionCreateReq
	if !response.BindCheckStruct(c, &req) {
		return
	}
	userId, ok := middleware.GetUserIDFromContext(c)
	if !ok {
		global.Log.Warn("用户操作失败：无法获取当前登录用户信息")
		response.FailWithMessage("无法获取当前登录用户信息", c)
		return
	}
	req.UserID = userId
	if err := userActionService.Create(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户操作成功", c)
}

func (u *UserActionApi) Delete(c *gin.Context) {
	actionId, ok := routex.QueryUint(c, "id")
	if !ok {
		return
	}
	userId, _ := middleware.GetUserIDFromContext(c)
	if err := userActionService.Delete(c.Request.Context(), actionId, userId); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户操作删除成功", c)
}

func (u *UserActionApi) List(c *gin.Context) {
	var req request.ActionListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("用户操作请求参数异常", zap.Error(err))
		response.FailWithMessage("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Error("用户操作请求参数异常", zap.String("errMsg", errMsg))
		response.FailWithMessage(errMsg, c)
		return
	}
	userId, _ := middleware.GetUserIDFromContext(c)
	list, total, err := userActionService.List(req, userId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(list, total, req.Page, req.PageSize, "用户操作列表获取成功", c)
}
