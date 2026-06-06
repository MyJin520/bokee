package base

import (
	"gin-admin/internal/mods/request"
	"gin-admin/internal/mods/response"
	"gin-admin/internal/network/middleware"
	"gin-admin/internal/network/service/base"
	"gin-admin/pkg/jwtx"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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

func (u *UserApi) ForgetPassword(c *gin.Context) {
	var req request.ForgetPasswordReq
	err := c.BindJSON(&req)
	if err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	cruId, _ := middleware.GetUserIDFromContext(c)
	err = userService.ForgetPassword(req, cruId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("密码重置成功", c)
}

// BindRoles 用户角色绑定
func (u *UserApi) BindRoles(c *gin.Context) {
	var req request.UserRoleBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}

	if err := userService.BindRoles(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("角色绑定成功", c)
}

func (u *UserApi) GetInfo(c *gin.Context) {
	var (
		userID uint
		err    error
		ok     bool
	)

	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		parsed, parseErr := strconv.ParseUint(userIDStr, 10, 32)
		if parseErr != nil {
			response.FailWithMessage("无效的用户ID", c)
			return
		}
		userID = uint(parsed)
	} else {
		userID, ok = middleware.GetUserIDFromContext(c)
		if !ok {
			response.FailWithMessage("无法获取当前用户信息", c)
			return
		}
	}

	user, err := userService.GetInfo(userID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(user, "获取用户信息成功", c)
}

func (u *UserApi) List(c *gin.Context) {
	var req request.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(http.StatusBadRequest, "请求参数异常>"+err.Error(), c)
		return
	}
	req.Normalize()

	users, total, err := userService.List(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithPage(users, total, req.Page, req.PageSize, "获取用户列表成功", c)
}
