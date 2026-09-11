package base

import (
	"bokee/global"
	"bokee/internal/mods/request"
	"bokee/internal/mods/response"
	"bokee/internal/network/middleware"
	"bokee/internal/network/service/base"
	"bokee/pkg/jwtx"
	"bokee/pkg/verifyx"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

type UserApi struct{}

var userService = &base.UserService{}

func (u *UserApi) Create(c *gin.Context) {
	var req request.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("用户创建请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("用户创建参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}
	if err := userService.Create(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户创建成功", c)
}

func (u *UserApi) Login(c *gin.Context) {
	var req request.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("用户登录请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("用户登录参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
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
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("Token 解析请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("Token 解析参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}
	claims, err := jwtx.ParseToken(req.Token)
	if err != nil {
		global.Log.Error("解析 Token 失败", zap.Error(err))
		response.FailWithMessage("无效的认证令牌", c)
		return
	}
	response.OkWithData(claims, "解析成功", c)
}

func (u *UserApi) Update(c *gin.Context) {
	var req request.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("用户更新请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("用户更新参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}
	userID, _ := middleware.GetUserIDFromContext(c)
	if err := userService.Update(req, userID); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户更新成功", c)
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
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("重置密码请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	errMsg := verifyx.CheckStruct(req)
	if errMsg != "" {
		global.Log.Warn("重置密码参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
		return
	}
	cruId, _ := middleware.GetUserIDFromContext(c)
	if err := userService.ForgetPassword(req, cruId); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("密码重置成功", c)
}

// BindRoles 用户角色绑定
func (u *UserApi) BindRoles(c *gin.Context) {
	var req request.UserRoleBindReq
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("角色绑定请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("角色绑定参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
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
			global.Log.Warn("用户ID格式错误", zap.String("userID", userIDStr))
			response.FailWithRequest("无效的用户ID", c)
			return
		}
		userID = uint(parsed)
	} else {
		userID, ok = middleware.GetUserIDFromContext(c)
		if !ok {
			global.Log.Warn("无法获取当前用户信息")
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
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("用户列表请求参数异常", zap.Error(err))
		response.FailWithRequest("请求参数异常", c)
		return
	}
	if errMsg := verifyx.CheckStruct(req); errMsg != "" {
		global.Log.Warn("用户列表参数校验失败", zap.String("err", errMsg))
		response.FailWithRequest(errMsg, c)
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
