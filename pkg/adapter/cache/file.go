package cache

import (
	"github.com/bhojpur/web/pkg/client/cache"
)

// NewFileCache Create new file cache with no config.
// the level and expiry need set in method StartAndGC as config string.
func NewFileCache() Cache {
	//    return &FileCache{CachePath:FileCachePath, FileSuffix:FileCacheFileSuffix}
	return CreateNewToOldCacheAdapter(cache.NewFileCache())
}

func init() {
	Register("file", NewFileCache)
}
