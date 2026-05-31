package redisx

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gin-admin/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	KeyPrefix     = "gin-admin"      // Key统一前缀
	EmptyCacheTTL = 30 * time.Second // 空数据缓存时间（防穿透）
	ErrPrefix     = "redis: "        // 错误统一前缀
)

// getClient 安全获取Redis客户端
func getClient() *redis.Client {
	if global.Redis == nil {
		panic("Redis 未初始化，请先调用 initRedis()")
	}
	return global.Redis
}

// minDuration 取最小时长
func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

// BuildKey 生成带命名空间的Key
func BuildKey(parts ...string) string {
	var validParts []string
	validParts = append(validParts, KeyPrefix)
	for _, p := range parts {
		if p != "" {
			validParts = append(validParts, p)
		}
	}
	return strings.Join(validParts, ":")
}

func Exists(ctx context.Context, key string) (bool, error) {
	n, err := getClient().Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf(ErrPrefix+"Exists failed: %w", err)
	}
	return n > 0, nil
}

func Get(ctx context.Context, key string) (string, error) {
	val, err := getClient().Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf(ErrPrefix+"Get failed: %w", err)
	}
	return val, nil
}

func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if err := getClient().Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf(ErrPrefix+"Set failed: %w", err)
	}
	return nil
}

func Delete(ctx context.Context, keys ...string) error {
	if err := getClient().Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf(ErrPrefix+"Del failed: %w", err)
	}
	return nil
}

func TTL(ctx context.Context, key string) (time.Duration, error) {
	d, err := getClient().TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf(ErrPrefix+"TTL failed: %w", err)
	}
	return d, nil
}

func Expire(ctx context.Context, key string, expiration time.Duration) error {
	if err := getClient().Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf(ErrPrefix+"Expire failed: %w", err)
	}
	return nil
}

func Incr(ctx context.Context, key string) (int64, error) {
	n, err := getClient().Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf(ErrPrefix+"Incr failed: %w", err)
	}
	return n, nil
}

func GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := getClient().Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return fmt.Errorf(ErrPrefix+"GetJSON failed: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf(ErrPrefix+"Unmarshal failed: %w", err)
	}
	return nil
}

func SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(ErrPrefix+"Marshal failed: %w", err)
	}
	if err := getClient().Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf(ErrPrefix+"SetJSON failed: %w", err)
	}
	return nil
}

// Remember 缓存穿透保护：缓存不存在则查询数据源并回填
func Remember(ctx context.Context, key string, expiration time.Duration, fallback func() (interface{}, error), dest interface{}) error {
	// 读取缓存
	if err := GetJSON(ctx, key, dest); err != nil {
		return err
	}
	// 命中缓存直接返回
	if !isEmptyValue(dest) {
		return nil
	}

	// 未命中，查询数据源
	data, err := fallback()
	if err != nil {
		return fmt.Errorf(ErrPrefix+"Remember fallback failed: %w", err)
	}

	// 回填缓存
	if data == nil {
		// 缓存空值，防止缓存穿透
		if err := SetJSON(ctx, key, struct{}{}, minDuration(expiration, EmptyCacheTTL)); err != nil {
			global.Log.Warn("缓存空值失败", zap.String("key", key), zap.Error(err))
		}
		return nil
	}

	// 缓存正常数据
	if err := SetJSON(ctx, key, data, expiration); err != nil {
		global.Log.Warn("回填缓存失败", zap.String("key", key), zap.Error(err))
	}

	// 赋值给目标变量
	_ = GetJSON(ctx, key, dest)
	return nil
}

// tokenBlacklistKey 生成 token 黑名单 key
func tokenBlacklistKey(tokenString string) string {
	h := sha256.Sum256([]byte(tokenString))
	return BuildKey("token", "blacklist", fmt.Sprintf("%x", h))
}

// BlacklistToken 将 token 加入黑名单，TTL 自动设为 token 剩余有效期
func BlacklistToken(ctx context.Context, tokenString string, ttl time.Duration) error {
	key := tokenBlacklistKey(tokenString)
	if err := getClient().Set(ctx, key, "1", ttl).Err(); err != nil {
		return fmt.Errorf(ErrPrefix+"BlacklistToken failed: %w", err)
	}
	return nil
}

// IsTokenBlacklisted 检查 token 是否已被拉黑
func IsTokenBlacklisted(ctx context.Context, tokenString string) (bool, error) {
	key := tokenBlacklistKey(tokenString)
	n, err := getClient().Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf(ErrPrefix+"IsTokenBlacklisted failed: %w", err)
	}
	return n > 0, nil
}

// isEmptyValue 判断是否为空结构体/空值
func isEmptyValue(v interface{}) bool {
	bytes, err := json.Marshal(v)
	if err != nil {
		return true
	}
	return string(bytes) == "{}"
}
