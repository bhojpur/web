// Package cache provide a Cache interface and some implement engine
// Usage:
//
// import(
//   "github.com/bhojpur/web/pkg/cache"
// )
//
// bm, err := cache.NewCache("memory", `{"interval":60}`)
//
// Use it like this:
//
//	bm.Put("bhojpur", 1, 10 * time.Second)
//	bm.Get("bhojpur")
//	bm.IsExist("bhojpur")
//	bm.Delete("bhojpur")
package cache

import (
	"context"
	"fmt"
	"time"
)

// Cache interface contains all behaviors for cache adapter.
// usage:
//	cache.Register("file",cache.NewFileCache) // this operation is run in init method of file.go.
//	c,err := cache.NewCache("file","{....}")
//	c.Put("key",value, 3600 * time.Second)
//	v := c.Get("key")
//
//	c.Incr("counter")  // now is 1
//	c.Incr("counter")  // now is 2
//	count := c.Get("counter").(int)
type Cache interface {
	// Get a cached value by key.
	Get(ctx context.Context, key string) (interface{}, error)
	// GetMulti is a batch version of Get.
	GetMulti(ctx context.Context, keys []string) ([]interface{}, error)
	// Set a cached value with key and expire time.
	Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error
	// Delete cached value by key.
	Delete(ctx context.Context, key string) error
	// Increment a cached int value by key, as a counter.
	Incr(ctx context.Context, key string) error
	// Decrement a cached int value by key, as a counter.
	Decr(ctx context.Context, key string) error
	// Check if a cached value exists or not.
	IsExist(ctx context.Context, key string) (bool, error)
	// Clear all cache.
	ClearAll(ctx context.Context) error
	// Start gc routine based on config string settings.
	StartAndGC(config string) error
}

// Instance is a function create a new Cache Instance
type Instance func() Cache

var adapters = make(map[string]Instance)

// Register makes a cache adapter available by the adapter name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

// NewCache creates a new cache driver by adapter name and config string.
// config: must be in JSON format such as {"interval":360}.
// Starts gc automatically.
func NewCache(adapterName, config string) (adapter Cache, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	err = adapter.StartAndGC(config)
	if err != nil {
		adapter = nil
	}
	return
}
