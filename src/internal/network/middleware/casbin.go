package middleware

import (
	"gin-admin/internal/mods/response"
	"gin-admin/pkg/casbinx"
	"net/http"

	"gin-admin/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 JWT 中间件注入的上下文中获取用户名
		username, ok := GetUsernameFromContext(c)
		if !ok {
			global.Log.Warn("Casbin 中间件：上下文中未找到用户名，拒绝访问")
			response.Fail(http.StatusForbidden, "未认证或令牌无效", c)
			c.Abort()
			return
		}

		// 2. 构造 Casbin 所需的三个参数：sub（主体）、obj（资源路径）、act（动作）
		sub := username
		obj := c.Request.URL.Path // 例如 /pri/user/edit
		act := c.Request.Method   // 例如 POST

		// 3. 执行权限校验
		allowed, err := casbinx.Enforce(sub, obj, act)
		if err != nil {
			global.Log.Error("Casbin 权限校验异常",
				zap.String("sub", sub),
				zap.String("obj", obj),
				zap.String("act", act),
				zap.Error(err),
			)
			response.Fail(http.StatusInternalServerError, "权限服务异常", c)
			c.Abort()
			return
		}

		if !allowed {
			global.Log.Warn("Casbin 拒绝访问",
				zap.String("sub", sub),
				zap.String("obj", obj),
				zap.String("act", act),
			)
			response.Fail(http.StatusForbidden, "无此操作权限", c)
			c.Abort()
			return
		}

		c.Next()
	}
}
