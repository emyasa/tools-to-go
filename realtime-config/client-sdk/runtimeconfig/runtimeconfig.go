package runtimeconfig

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	client *redis.Client
	channel string
}

func NewRedisOptions(
	client *redis.Client,
	channel string,
) (*RedisOptions, error) {
	if client == nil {
		return nil, errors.New("client is required")
	}

	if channel == "" {
		return nil, errors.New("channel is required")
	}

	return &RedisOptions{
		client: client,
		channel: channel,
	}, nil
}

type Loader[T any] func(context.Context, *redis.Client) (*T, error)

type Store[T any] struct {
	current atomic.Pointer[T]
	loader Loader[T]
}

func (s *Store[T]) watch(
	ctx context.Context,
	redisOptions *RedisOptions,
) {
	sub := redisOptions.client.Subscribe(ctx, redisOptions.channel)
	defer sub.Close()

	for range sub.Channel() {
		cfg, err := s.loader(ctx, redisOptions.client)
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
	redisOptions *RedisOptions,
	loader Loader[T],
) (*Store[T], error) {
	if redisOptions == nil {
		return nil, errors.New("redisOptions is required")
	}

	cfg, err := loader(ctx, redisOptions.client)
	if err != nil {
		return nil, err
	}

	s := &Store[T] {
		loader: loader,
	}

	s.current.Store(cfg)
	go s.watch(ctx, redisOptions)

	return s, nil
}
