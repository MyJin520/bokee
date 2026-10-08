package cachex

import (
	"context"
	"fmt"

	"bokee/global"
	"bokee/pkg/redisx"

	"go.uber.org/zap"
)

// 业务缓存命名空间（Key 形如 bokee:<namespace>:<biz>:<id>），新模块在此登记
const (
	NSUser    = "user"
	NSArticle = "article"
	NSRole    = "role"
)

// Key 生成业务缓存 Key（自动带 bokee 前缀，空段自动跳过）。
// 段值支持 string 与整数（如 uint id），整数自动转为十进制字符串。
func Key(parts ...any) string {
	strs := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == nil {
			continue
		}
		strs = append(strs, fmt.Sprint(p))
	}
	return redisx.BuildKey(strs...)
}

// Invalidate 批量失效缓存 Key（先写库、后删缓存，保证最终一致；失败仅记日志不阻断主流程）
func Invalidate(ctx context.Context, keys ...string) {
	if len(keys) == 0 {
		return
	}
	if err := redisx.Delete(ctx, keys...); err != nil {
		global.Log.Error("失效缓存失败", zap.Strings("keys", keys), zap.Error(err))
	}
}
