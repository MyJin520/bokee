package base

import (
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/internal/network/middleware"
	"gin-admin/internal/network/service/base"
	"gin-admin/pkg/jwtx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserApi struct{}

var userService = &base.UserService{}

func (u *UserApi) Register(c *gin.Context) {
	var req request.UserRegisterReq
	err := c.BindJSON(&req)
	if err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	err = userService.Register(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户注册成功", c)
}

func (u *UserApi) Login(c *gin.Context) {
	var req request.UserLoginReq
	err := c.BindJSON(&req)
	if err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	jwtResponse, err := userService.Login(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(jwtResponse, "登录成功", c)
}

// ParseToken TODO 测试解析Token
func (u *UserApi) ParseToken(c *gin.Context) {
	var req request.TokenParsingReq
	err := c.BindJSON(&req)
	if err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(claims, "解析成功", c)
}

func (u *UserApi) Edit(c *gin.Context) {
	var req request.UserEditReq
	err := c.BindJSON(&req)
	if err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	userID, _ := middleware.GetUserIDFromContext(c)
	err = userService.Edit(req, userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户编辑成功", c)
}

func (u *UserApi) Logout(c *gin.Context) {
	err := userService.Logout(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("登出成功", c)
}

// GetAllPriRoles todo 仅管理员访问能访问
func (u *UserApi) GetAllPriRoles(c *gin.Context) {
	var page request.PageReq
	if err := c.ShouldBindQuery(&page); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	page.Normalize()

	roles, total, err := userService.GetAllPriRoles(page)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(roles, total, page.Page, page.PageSize, "获取角色列表成功", c)
}
