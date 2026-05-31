package middleware

import (
	"crypto/sha256"
	"fmt"
	"gin-admin/global"
	"gin-admin/internal/mods/response"
	"gin-admin/pkg/jwtx"
	"go.uber.org/zap"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// blacklistKey 生成 token 黑名单 key（使用 SHA256 哈希，不依赖 jti）
// 这样无论旧 token（无 jti）还是新 token 都能被正确拉黑
func blacklistKey(tokenString string) string {
	h := sha256.Sum256([]byte(tokenString))
	return "gin-admin:token:blacklist:" + fmt.Sprintf("%x", h)
}

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			global.Log.Warn("请求未携带 Authorization 头")
			response.Fail(http.StatusUnauthorized, "未提供认证令牌", c)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			global.Log.Warn("Authorization 头格式错误", zap.String("header", authHeader))
			response.Fail(http.StatusUnauthorized, "认证令牌格式错误", c)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析并验证 token
		claims, err := jwtx.ParseToken(tokenString)
		if err != nil {
			global.Log.Warn("无效的 JWT 令牌", zap.Error(err))
			response.Fail(http.StatusUnauthorized, "无效或过期的认证令牌", c)
			c.Abort()
			return
		}

		// 检查 Redis 黑名单（token 是否已登出）
		if global.Redis != nil {
			key := blacklistKey(tokenString)
			exists, err := global.Redis.Exists(c, key).Result()
			if err != nil {
				global.Log.Warn("Redis 黑名单查询失败", zap.String("jti", claims.ID), zap.Error(err))
				// 查询失败放行，不阻断请求
			} else if exists > 0 {
				global.Log.Warn("Token 已被登出", zap.String("jti", claims.ID), zap.Uint("userID", claims.UserID))
				response.Fail(http.StatusUnauthorized, "令牌已失效，请重新登录", c)
				c.Abort()
				return
			}
		}

		// 将用户信息存入 Gin 上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// GetUserIDFromContext 从上下文获取当前登录用户的 ID
func GetUserIDFromContext(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUsernameFromContext 从上下文获取当前登录用户名
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}
