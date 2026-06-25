package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emyasa/tools-to-go/realtime-config/client-sdk/runtimeconfig"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type NewServiceConfig struct {
	Template string `yaml:"template"`
}

type GeneralConfig struct {
	Key string `yaml:"key"`
}

type Config struct {
	NewService NewServiceConfig
	General GeneralConfig
}

func load(ctx context.Context, rdb *redis.Client) (*Config, error) {
	newService, _ := rdb.Get(ctx, "new-service/config.yaml").Result()
	general, _ := rdb.Get(ctx, "general.yaml").Result()

	var cfg Config
	yaml.Unmarshal([]byte(newService), &cfg.NewService)
	yaml.Unmarshal([]byte(general), &cfg.General)

	return &cfg, nil
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rdb.Close()

	opts, err := runtimeconfig.NewRedisOptions(rdb, "client-example")
	if err != nil {
		panic(err)
	}

	store, err := runtimeconfig.Watch(ctx, opts, load)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("currentConfig: %+v\n", store.Get())
	})

	fmt.Println("Listening on port: 8080")
	http.ListenAndServe(":8080", nil)
}
