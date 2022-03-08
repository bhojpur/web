package app

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"sync"
	"time"

	"github.com/bhojpur/web/pkg/app/cache"
)

// PreRenderCache is the interface that describes a cache that stores
// pre-rendered resources.
type PreRenderCache interface {
	// Get returns the item at the given path.
	Get(ctx context.Context, path string) (PreRenderedItem, bool)

	// Set stored the item at the given path.
	Set(ctx context.Context, i PreRenderedItem)
}

// PreRenderedItem represent an item that is stored in a PreRenderCache.
type PreRenderedItem struct {
	// The request path.
	Path string

	// The response content type.
	ContentType string

	// The response content encoding.
	ContentEncoding string

	// The response body.
	Body []byte
}

// Len return the body length.
func (r PreRenderedItem) Size() int {
	return len(r.Body)
}

// NewPreRenderLRUCache creates an in memory LRU cache that stores items for the
// given duration. If provided, on eviction functions are called when item are
// evicted.
func NewPreRenderLRUCache(size int, itemTTL time.Duration, onEvict ...func(path string, i PreRenderedItem)) PreRenderCache {
	return &preRenderLRUCache{
		LRU: cache.LRU{
			MaxSize: size,
			ItemTTL: itemTTL,
			OnEvict: func(path string, i cache.Item) {
				item := i.(PreRenderedItem)
				for _, fn := range onEvict {
					fn(path, item)
				}
			},
		},
	}

}

type preRenderLRUCache struct {
	cache.LRU
}

func (c *preRenderLRUCache) Get(ctx context.Context, path string) (PreRenderedItem, bool) {
	i, ok := c.LRU.Get(ctx, path)
	if !ok {
		return PreRenderedItem{}, false
	}
	return i.(PreRenderedItem), true
}

func (c *preRenderLRUCache) Set(ctx context.Context, i PreRenderedItem) {
	c.LRU.Set(ctx, i.Path, i)
}

type preRenderCache struct {
	mu    sync.RWMutex
	items map[string]PreRenderedItem
}

func newPreRenderCache(size int) *preRenderCache {
	return &preRenderCache{
		items: make(map[string]PreRenderedItem, size),
	}
}
func (c *preRenderCache) Set(ctx context.Context, i PreRenderedItem) {
	c.mu.Lock()
	c.items[i.Path] = i
	c.mu.Unlock()
}
func (c *preRenderCache) Get(ctx context.Context, path string) (PreRenderedItem, bool) {
	c.mu.Lock()
	i, ok := c.items[path]
	c.mu.Unlock()
	return i, ok
}
