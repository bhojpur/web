package cache

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
	"sort"
	"sync"
	"time"
)

// LRU represents an LRU (Least Recently Used) cache implementation.
type LRU struct {
	// The maximum cache size in bytes. Default is 16MB.
	MaxSize int

	// The duration while an item is cached.
	ItemTTL time.Duration

	// The function called when an item is evicted.
	OnEvict func(key string, i Item)

	once     sync.Once
	mutex    sync.Mutex
	size     int
	items    map[string]*lruItem
	priority []*lruItem
}

func (c *LRU) Get(ctx context.Context, key string) (Item, bool) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	i, isCached := c.items[key]
	if !isCached || i.expiresAt.Before(time.Now()) {
		return nil, false
	}

	i.count++
	return i.value, true
}

func (c *LRU) Set(ctx context.Context, key string, i Item) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		i.expiresAt = time.Now()
	}

	if c.size+i.Size() > c.MaxSize {
		c.free(i.Size())
	}

	c.add(&lruItem{
		key:       key,
		count:     1,
		expiresAt: time.Now().Add(c.ItemTTL),
		value:     i,
	})
}

func (c *LRU) Del(ctx context.Context, key string) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		i.expiresAt = time.Now()
		c.free(i.value.Size())
	}
}

func (c *LRU) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.items)
}

func (c *LRU) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.size
}

func (c *LRU) init() {
	if c.MaxSize <= 0 {
		c.MaxSize = 16000000
	}

	c.items = make(map[string]*lruItem, 64)
	c.priority = make([]*lruItem, 0, len(c.items))
}

func (c *LRU) free(size int) {
	now := time.Now()
	sortLRUItems(now, c.priority)

	for len(c.priority) != 0 {
		lastItem := c.priority[len(c.priority)-1]
		if lastItem.IsExpired(now) {
			c.removeLastItem()
			continue
		}

		if c.size+size <= c.MaxSize {
			return
		}
		c.removeLastItem()
		if c.OnEvict != nil {
			c.OnEvict(lastItem.key, lastItem.value)
		}
	}
}

func (c *LRU) removeLastItem() {
	i := len(c.priority) - 1
	item := c.priority[i]
	c.priority[i] = nil
	c.priority = c.priority[:i]
	delete(c.items, item.key)
	c.size -= item.value.Size()
}

func (c *LRU) add(i *lruItem) {
	c.items[i.key] = i
	c.priority = append(c.priority, i)
	c.size += i.value.Size()
}

type lruItem struct {
	key       string
	count     int
	expiresAt time.Time
	value     Item
}

func (i *lruItem) priority(now time.Time) int {
	if i.IsExpired(now) {
		return 0
	}
	return i.count
}

func (i *lruItem) IsExpired(now time.Time) bool {
	return i.expiresAt.Before(now)
}

func sortLRUItems(now time.Time, v []*lruItem) {
	sort.Slice(v, func(a, b int) bool {
		return v[a].priority(now) > v[b].priority(now)
	})
}
