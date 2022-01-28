// depend on github.com/gomodule/redigo/redis
//
// go install github.com/gomodule/redigo/redis
//
// Usage:
// import(
//   _ "github.com/bhojpur/web/pkg/cache/redis"
//   "github.com/bhojpur/web/pkg/cache"
// )
//
//  bm, err := cache.NewCache("redis", `{"conn":"127.0.0.1:11211"}`)
package redis

import (
	"github.com/bhojpur/web/pkg/adapter/cache"
	redis2 "github.com/bhojpur/web/pkg/client/cache/redis"
)

var (
	// DefaultKey the collection name of redis for cache adapter.
	DefaultKey = "bcacheRedis"
)

// NewRedisCache create new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return cache.CreateNewToOldCacheAdapter(redis2.NewRedisCache())
}

func init() {
	cache.Register("redis", NewRedisCache)
}
