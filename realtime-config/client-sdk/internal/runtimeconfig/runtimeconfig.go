package runtimeconfig

import (
	"context"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

type Options struct {
	RedisAddr string
	Channel string
}

func NewOptions(redisAddr, channel string) Options {
	if redisAddr == "" {
		panic("redisAddr is required")
	}

	if channel == "" {
		panic("channel is required")
	}

	return Options{
		RedisAddr: redisAddr,
		Channel: channel,
	}
}

type Loader[T any] func(context.Context, *redis.Client) (*T, error)

type Store[T any] struct {
	current atomic.Pointer[T]
	loader Loader[T]
}

func (s *Store[T]) watch(
	ctx context.Context,
	rdb *redis.Client,
	channel string,
) {
	sub := rdb.Subscribe(ctx, channel)
	defer sub.Close()

	for range sub.Channel() {
		cfg, err := s.loader(ctx, rdb)
		if err != nil {
			continue
		}

		s.current.Store(cfg)
	}
}

func (s *Store[T]) Get() (*T) {
	return s.current.Load()
}

func Watch[T any](
	ctx context.Context,
	options Options,
	loader Loader[T],
) (*Store[T], error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: options.RedisAddr,
	})

	cfg, err := loader(ctx, rdb)
	if err != nil {
		return nil, err
	}

	s := &Store[T] {
		loader: loader,
	}

	s.current.Store(cfg)
	go s.watch(ctx, rdb, options.Channel)

	return s, nil
}
