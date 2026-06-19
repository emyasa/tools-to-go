# TODO
[x] Sync Repo
[x] Load Repo Files into a key-value map, where key is the directory + filename and value is the string content
[x] Spin up a redis container: `docker run -d --name redis -p 6379:6379 redis:7-alpine`
[x] Set up redis client
[x] Load the files key-value map into redis
[x] Verify that the map has been loaded to redis
```
redis-cli -h localhost -p 6379
KEYS config-server:*
```
