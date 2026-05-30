package base

import (
	"gin-admin/internal/mods/request"
	"gin-admin/internal/network/service/base"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserApi struct{}

var userService = &base.UserService{}

func (u *UserApi) Register(c *gin.Context) {
	var req request.UserRegisterOrEditReq
	err := c.BindJSON(&req)
	if err != nil {
		return
	}
	_ = userService.Register(req)
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}
