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
