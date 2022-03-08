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
	"sync"
	"time"
)

type Expire struct {
	// The duration while an item is cached.
	ItemTTL time.Duration

	once  sync.Once
	mutex sync.RWMutex
	size  int
	items map[string]*memItem
	queue []*memItem
}

func (c *Expire) Get(ctx context.Context, key string) (Item, bool) {
	c.once.Do(c.init)
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	i, isCached := c.items[key]
	if !isCached || i.expiresAt.Before(time.Now()) {
		return nil, false
	}
	return i.value, true
}

func (c *Expire) Set(ctx context.Context, key string, i Item) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		c.del(i)
	}

	c.expire()
	c.add(&memItem{
		key:       key,
		expiresAt: time.Now().Add(c.ItemTTL),
		value:     i,
	})
}

func (c *Expire) Del(ctx context.Context, key string) {
	c.once.Do(c.init)
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if i, isCached := c.items[key]; isCached {
		c.del(i)
		c.expire()
	}
}

func (c *Expire) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.items)
}

func (c *Expire) Size() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.size
}

func (c *Expire) init() {
	c.items = make(map[string]*memItem)
}

func (c *Expire) add(i *memItem) {
	c.items[i.key] = i
	c.queue = append(c.queue, i)
	c.size += i.value.Size()
}

func (c *Expire) del(i *memItem) {
	if i.value != nil {
		delete(c.items, i.key)
		c.size -= i.value.Size()
		i.value = nil
	}
}

func (c *Expire) expire() {
	now := time.Now()

	i := 0
	for i < len(c.queue) && c.queue[i].isExpired(now) {
		c.del(c.queue[i])
		i++
	}

	copy(c.queue, c.queue[i:])
	c.queue = c.queue[:len(c.queue)-i]
}

type memItem struct {
	key       string
	expiresAt time.Time
	value     Item
}

func (i *memItem) isExpired(now time.Time) bool {
	return i.value == nil || i.expiresAt.Before(now)
}
