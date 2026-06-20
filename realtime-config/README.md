# realtime-config

Realtime-config aims to aid in achieving a familiar way of managing non-sensitive configurations in a Git repo and loading them without a server restart.

## How it works
```
┌─────────────────┐     ┌──────────────────┐     ┌────────────────┐
│  config-loader  │ ──▶ │      Redis       │ ──▶ │  client-sdk    │
│  (Git → Redis)  │     │  (key/value +    │     │  (Redis → app) │
│                 │     │   pub/sub)       │     │                │
└─────────────────┘     └──────────────────┘     └────────────────┘
```

## Usage

### 1. config-loader - push config from Git to Redis (see loader_test.go)
```go
gitOpts, _ := loader.NewGitOptions("/path/to/id_ed25519", "git@github.com:user/config-repo", "main")
redisOpts, _ := loader.NewRedisOptions(rdb, "prefix/", []string{"target-client-channel"})
loader.LoadFromGit(context.Background(), gitOpts, redisOpts)
```

### 2. client-sdk - watch for changes in your app (see examples/client/main.go)
```go
type Config struct {
    ...
}

func load(ctx context.Context, rdb *redis.Client) (*Config, error) {
	newService, _ := rdb.Get(ctx, "prefix/new-service/config.yaml").Result()
	general, _ := rdb.Get(ctx, "prefix/general.yaml").Result()

	var cfg Config
	yaml.Unmarshal([]byte(newService), &cfg.NewService)
	yaml.Unmarshal([]byte(general), &cfg.General)

	return &cfg, nil
}

opts, _ := runtimeconfig.NewRedisOptions(rdb, "target-client-channel")
store, _ := runtimeconfig.Watch(ctx, opts, load)
```
