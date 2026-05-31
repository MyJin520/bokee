package initialize

import (
	"context"
	"time"

	"gin-admin/config"
	"gin-admin/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// initRedis 初始化 Redis 客户端连接
func initRedis(cfg *config.Config) {
	opt := &redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  time.Duration(cfg.Redis.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Redis.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Redis.WriteTimeout) * time.Second,
	}

	client := redis.NewClient(opt)

	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		global.Log.Warn("Redis 连接失败，将使用降级模式运行", zap.String("addr", cfg.Redis.Addr), zap.Error(err))
		global.Redis = nil
		return
	}

	global.Redis = client
	global.Log.Info("Redis 连接成功", zap.String("addr", cfg.Redis.Addr), zap.Int("db", cfg.Redis.DB))
}
