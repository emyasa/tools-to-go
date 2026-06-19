module github.com/emyasa/tools-to-go/realtime-config/examples/client

go 1.26.4

require (
	github.com/emyasa/tools-to-go/realtime-config/client-sdk v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.20.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

replace github.com/emyasa/tools-to-go/realtime-config/client-sdk => ../../client-sdk
