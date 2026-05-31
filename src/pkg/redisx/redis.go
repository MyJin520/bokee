package redisx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gin-admin/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// getClient 安全获取 Redis 客户端，未初始化时 panic
func getClient() *redis.Client {
	if global.Redis == nil {
		panic("Redis 未初始化: global.Redis 为 nil，请确认 initRedis() 已调用且连接成功")
	}
	return global.Redis
}

// Exists 检查 key 是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	n, err := getClient().Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis Exists 失败: %w", err)
	}
	return n > 0, nil
}

// Get 获取字符串值
func Get(ctx context.Context, key string) (string, error) {
	val, err := getClient().Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("redis Get 失败: %w", err)
	}
	return val, nil
}

// Set 设置字符串值
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := getClient().Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis Set 失败: %w", err)
	}
	return nil
}

// Delete 删除一个或多个 key
func Delete(ctx context.Context, keys ...string) error {
	err := getClient().Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("redis Del 失败: %w", err)
	}
	return nil
}

// TTL 获取 key 剩余生存时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	d, err := getClient().TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis TTL 失败: %w", err)
	}
	return d, nil
}

// Expire 设置 key 的过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := getClient().Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis Expire 失败: %w", err)
	}
	return nil
}

// Incr 自增（原子操作）
func Incr(ctx context.Context, key string) (int64, error) {
	n, err := getClient().Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis Incr 失败: %w", err)
	}
	return n, nil
}

// GetJSON 从 Redis 获取 JSON 字符串并反序列化到 dest
// dest 必须是指针类型
func GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := getClient().Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("redis GetJSON 失败: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("redis GetJSON 反序列化失败: %w", err)
	}
	return nil
}

// SetJSON 将 value 序列化为 JSON 并存入 Redis
func SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redis SetJSON 序列化失败: %w", err)
	}
	if err := getClient().Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("redis SetJSON 失败: %w", err)
	}
	return nil
}

// Remember 缓存穿透保护
// 1. 优先查 Redis（key 反序列化到 dest）
// 2. 未命中则调用 fallback 获取数据，回填缓存后返回
// 3. fallback 返回 nil, nil 表示数据不存在（缓存空值，防止缓存穿透，空值过期时间缩短）
//
// 用法:
//
//	var user User
//	err := redisx.Remember(ctx, "user:1", 10*time.Minute, func() (interface{}, error) {
//	    return db.FindUser(1)
//	}, &user)
func Remember(ctx context.Context, key string, expiration time.Duration, fallback func() (interface{}, error), dest interface{}) error {
	// 1. 尝试从缓存读取
	found, err := getFromCache(ctx, key, dest)
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	// 2. 未命中 → 调用 fallback
	data, err := fallback()
	if err != nil {
		return fmt.Errorf("redis Remember fallback 失败: %w", err)
	}

	// 3. 回填缓存
	if data == nil {
		// 数据不存在 → 缓存空值（短过期，防穿透）
		cacheMissErr := SetJSON(ctx, key, struct{}{}, min(expiration, 30*time.Second))
		if cacheMissErr != nil {
			global.Log.Warn("Redis Remember 缓存空值失败", zap.String("key", key), zap.Error(cacheMissErr))
		}
		return nil
	}

	if err := SetJSON(ctx, key, data, expiration); err != nil {
		global.Log.Warn("Redis Remember 回填缓存失败", zap.String("key", key), zap.Error(err))
	}

	// 4. 反序列化到 dest 返回
	raw, _ := json.Marshal(data)
	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("redis Remember 反序列化失败: %w", err)
	}
	return nil
}

// getFromCache 尝试从缓存读取，返回是否命中
func getFromCache(ctx context.Context, key string, dest interface{}) (bool, error) {
	data, err := getClient().Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("redis Remember 读取缓存失败: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("redis Remember 反序列化缓存失败: %w", err)
	}
	return true, nil
}

// Lock 分布式锁
type Lock struct {
	key        string
	token      string // 随机 token，防止误释放
	expiration time.Duration
	client     *redis.Client
}

// NewLock 创建分布式锁
// defaultExpiration 锁的默认持有时间（业务超时自动释放，防止死锁）
func NewLock(key string, defaultExpiration time.Duration) *Lock {
	return &Lock{
		key:        key,
		token:      fmt.Sprintf("%d", time.Now().UnixNano()),
		expiration: defaultExpiration,
		client:     getClient(),
	}
}

// Acquire 尝试获取锁，阻塞等待
// 使用指数退避重试，最长等待 waitTimeout
func (l *Lock) Acquire(ctx context.Context, waitTimeout time.Duration) error {
	deadline := time.Now().Add(waitTimeout)
	backoff := 50 * time.Millisecond
	maxBackoff := 2 * time.Second

	for {
		ok, err := l.client.SetNX(ctx, l.key, l.token, l.expiration).Result()
		if err != nil {
			return fmt.Errorf("redis Lock Acquire 失败: %w", err)
		}
		if ok {
			return nil
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("redis Lock Acquire 超时: 等待 %v 后仍未能获取锁 %s", waitTimeout, l.key)
		}

		time.Sleep(backoff)
		if backoff < maxBackoff {
			backoff *= 2
		}
	}
}

// AcquireOnce 尝试获取锁一次，不阻塞
// 成功返回 true，失败返回 false（不报错）
func (l *Lock) AcquireOnce(ctx context.Context) (bool, error) {
	ok, err := l.client.SetNX(ctx, l.key, l.token, l.expiration).Result()
	if err != nil {
		return false, fmt.Errorf("redis Lock AcquireOnce 失败: %w", err)
	}
	return ok, nil
}

// Release 释放锁（使用 Lua 脚本，保证原子性 — 只释放自己的锁）
func (l *Lock) Release(ctx context.Context) error {
	// KEYS[1] = lock key, ARGV[1] = token
	script := redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("DEL", KEYS[1])
		else
			return 0
		end
	`)
	err := script.Run(ctx, l.client, []string{l.key}, l.token).Err()
	if err != nil {
		return fmt.Errorf("redis Lock Release 失败: %w", err)
	}
	return nil
}

// BuildKey 生成带命名空间的 Redis key
// 用法: redisx.BuildKey("user", fmt.Sprintf("%d", userID)) → "gin-admin:user:123"
func BuildKey(parts ...string) string {
	const prefix = "gin-admin"
	segments := append([]string{prefix}, parts...)
	// 计算总长度预分配
	total := 0
	for _, s := range segments {
		if s != "" {
			total += len(s) + 1
		}
	}
	buf := make([]byte, 0, total)
	for _, s := range segments {
		if s != "" {
			buf = append(buf, s...)
			buf = append(buf, ':')
		}
	}
	if len(buf) > 0 {
		buf = buf[:len(buf)-1] // 去掉末尾 ':'
	}
	return string(buf)
}
