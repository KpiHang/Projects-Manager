package dao

import (
	"context"
	"github.com/go-redis/redis/v8"
	"test.com/project-user/config"
	"time"
)

var Rc *RedisCache // 用Redis的方式实现Cache接口；

// RedisCache 依赖rdb(和redis的连接)
type RedisCache struct {
	rdb *redis.Client
}

func init() { // 使用redis先连接，就放到init中；
	rdb := redis.NewClient(config.Conf.GetRedisConfig())

	Rc = &RedisCache{
		rdb: rdb,
	}
}

func (rc *RedisCache) Put(ctx context.Context, key, value string, expire time.Duration) error {
	err := rc.rdb.Set(ctx, key, value, expire).Err()
	return err
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := rc.rdb.Get(ctx, key).Result()
	return result, err
}
