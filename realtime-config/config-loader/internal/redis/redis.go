package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func New(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

func LoadMapAsKeys(ctx context.Context, rdb *redis.Client, prefix string, data map[string]string) error {
	pipe := rdb.Pipeline()

	for k, v := range data {
		pipe.Set(ctx, prefix+":"+k, v, 0)
	}

	_, err := pipe.Exec(ctx)
	return err
}
