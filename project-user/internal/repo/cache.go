package repo

import (
	"context"
	"time"
)

type Cache interface {
	Put(ctx context.Context, key, value string, expire time.Duration) error // 存
	Get(ctx context.Context, key string) (string, error)                    // 取
}
