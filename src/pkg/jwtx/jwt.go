package jwtx

import (
	"errors"
	"gin-admin/global"
	"gin-admin/pkg/timex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义 Claims，可根据业务扩展
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

// getJwtSecret 安全获取 JWT 密钥
func getJwtSecret() []byte {
	if global.Config == nil {
		panic("JWT initialization failed: global.Config is nil, please ensure configuration is loaded before using jwtx package")
	}
	if global.Config.JWT.Secret == "" {
		panic("JWT initialization failed: secret is empty, please set a valid JWT secret in configuration")
	}
	return []byte(global.Config.JWT.Secret)
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string) (string, error) {
	expire, err := timex.ParseDuration(global.Config.JWT.Expire)
	if err != nil {
		return "", err
	}

	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Config.JWT.Issuer,
			Subject:   username + "kim",
			Audience:  global.Config.JWT.Audience,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJwtSecret())
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJwtSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken 刷新 token
func RefreshToken(oldTokenString string, extendDuration time.Duration) (string, error) {
	claims, err := ParseToken(oldTokenString)
	if err != nil {
		return "", errors.New("invalid old token, can't refresh")
	}

	// 生成新的过期时间
	newExpire := time.Now().Add(extendDuration)
	claims.ExpiresAt = jwt.NewNumericDate(newExpire)
	claims.IssuedAt = jwt.NewNumericDate(time.Now())

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString(getJwtSecret())
}

// VerifyToken 仅验证有效性，返回布尔值
func VerifyToken(tokenString string) bool {
	_, err := ParseToken(tokenString)
	return err == nil
}
