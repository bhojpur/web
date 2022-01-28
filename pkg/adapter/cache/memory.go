package cache

import (
	"github.com/bhojpur/web/pkg/client/cache"
)

// NewMemoryCache returns a new MemoryCache.
func NewMemoryCache() Cache {
	return CreateNewToOldCacheAdapter(cache.NewMemoryCache())
}

func init() {
	Register("memory", NewMemoryCache)
}
