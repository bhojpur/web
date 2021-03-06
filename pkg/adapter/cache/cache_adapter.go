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
	"time"

	"github.com/bhojpur/web/pkg/client/cache"
)

type newToOldCacheAdapter struct {
	delegate cache.Cache
}

func (c *newToOldCacheAdapter) Get(key string) interface{} {
	res, _ := c.delegate.Get(context.Background(), key)
	return res
}

func (c *newToOldCacheAdapter) GetMulti(keys []string) []interface{} {
	res, _ := c.delegate.GetMulti(context.Background(), keys)
	return res
}

func (c *newToOldCacheAdapter) Put(key string, val interface{}, timeout time.Duration) error {
	return c.delegate.Put(context.Background(), key, val, timeout)
}

func (c *newToOldCacheAdapter) Delete(key string) error {
	return c.delegate.Delete(context.Background(), key)
}

func (c *newToOldCacheAdapter) Incr(key string) error {
	return c.delegate.Incr(context.Background(), key)
}

func (c *newToOldCacheAdapter) Decr(key string) error {
	return c.delegate.Decr(context.Background(), key)
}

func (c *newToOldCacheAdapter) IsExist(key string) bool {
	res, err := c.delegate.IsExist(context.Background(), key)
	return res && err == nil
}

func (c *newToOldCacheAdapter) ClearAll() error {
	return c.delegate.ClearAll(context.Background())
}

func (c *newToOldCacheAdapter) StartAndGC(config string) error {
	return c.delegate.StartAndGC(config)
}

func CreateNewToOldCacheAdapter(delegate cache.Cache) Cache {
	return &newToOldCacheAdapter{
		delegate: delegate,
	}
}

type oldToNewCacheAdapter struct {
	old Cache
}

func (o *oldToNewCacheAdapter) Get(ctx context.Context, key string) (interface{}, error) {
	return o.old.Get(key), nil
}

func (o *oldToNewCacheAdapter) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	return o.old.GetMulti(keys), nil
}

func (o *oldToNewCacheAdapter) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return o.old.Put(key, val, timeout)
}

func (o *oldToNewCacheAdapter) Delete(ctx context.Context, key string) error {
	return o.old.Delete(key)
}

func (o *oldToNewCacheAdapter) Incr(ctx context.Context, key string) error {
	return o.old.Incr(key)
}

func (o *oldToNewCacheAdapter) Decr(ctx context.Context, key string) error {
	return o.old.Decr(key)
}

func (o *oldToNewCacheAdapter) IsExist(ctx context.Context, key string) (bool, error) {
	return o.old.IsExist(key), nil
}

func (o *oldToNewCacheAdapter) ClearAll(ctx context.Context) error {
	return o.old.ClearAll()
}

func (o *oldToNewCacheAdapter) StartAndGC(config string) error {
	return o.old.StartAndGC(config)
}

func CreateOldToNewAdapter(old Cache) cache.Cache {
	return &oldToNewCacheAdapter{
		old: old,
	}
}
