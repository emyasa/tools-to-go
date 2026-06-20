package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emyasa/tools-to-go/realtime-config/client-sdk/runtimeconfig"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type GeneralConfig struct {
	Sample string `yaml:"sample"`
}

type Config struct {
	General GeneralConfig
}

func load(ctx context.Context, rdb *redis.Client) (*Config, error) {
	general, _ := rdb.Get(ctx, "general.yaml").Result()

	var cfg Config
	yaml.Unmarshal([]byte(general), &cfg.General)

	return &cfg, nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	opts, err := runtimeconfig.NewRedisOptions(rdb, "general")
	if err != nil {
		panic(err)
	}

	cfg, err := runtimeconfig.Watch(ctx, opts, load)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		general := cfg.Get().General
		fmt.Printf("%+v\n", general)
	})

	fmt.Println("Listening on port: 8080")
	http.ListenAndServe(":8080", nil)
}
