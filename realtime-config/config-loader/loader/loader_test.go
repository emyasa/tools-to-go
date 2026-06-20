package loader

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestLoadConfigFromGit(t *testing.T) {
	gitOptions, err := NewGitOptions(
		"/Users/emyasa/.ssh/id_ed25519",
		"git@github.com:emyasa/scratch-config",
		"main",
	)

	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	redisOptions, err := NewRedisOptions(rdb, "", []string{"general"})

	if err != nil {
		panic(err)
	}

	LoadFromGit(
		context.Background(),
		gitOptions,
		redisOptions,
	)
}
