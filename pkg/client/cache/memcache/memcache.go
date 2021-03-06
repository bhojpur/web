package memcache

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

// depend on github.com/bradfitz/gomemcache/memcache
//
// go install github.com/bradfitz/gomemcache/memcache
//
// Usage:
// import(
//   _ "github.com/bhojpur/web/cache/memcache"
//   "github.com/bhojpur/web/cache"
// )
//
//  bm, err := cache.NewCache("memcache", `{"conn":"127.0.0.1:11211"}`)

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/bhojpur/web/pkg/client/cache"
)

// Cache Memcache adapter.
type Cache struct {
	conn     *memcache.Client
	conninfo []string
}

// NewMemCache creates a new memcache adapter.
func NewMemCache() cache.Cache {
	return &Cache{}
}

// Get get value from memcache.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return nil, err
		}
	}
	if item, err := rc.conn.Get(key); err == nil {
		return item.Value, nil
	} else {
		return nil, err
	}
}

// GetMulti gets a value from a key in memcache.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	rv := make([]interface{}, len(keys))
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return rv, err
		}
	}

	mv, err := rc.conn.GetMulti(keys)
	if err != nil {
		return rv, err
	}

	keysErr := make([]string, 0)
	for i, ki := range keys {
		if _, ok := mv[ki]; !ok {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, "the key isn't exist"))
			continue
		}
		rv[i] = mv[ki].Value
	}

	if len(keysErr) == 0 {
		return rv, nil
	}
	return rv, fmt.Errorf(strings.Join(keysErr, "; "))
}

// Put puts a value into memcache.
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	item := memcache.Item{Key: key, Expiration: int32(timeout / time.Second)}
	if v, ok := val.([]byte); ok {
		item.Value = v
	} else if str, ok := val.(string); ok {
		item.Value = []byte(str)
	} else {
		return errors.New("val only support string and []byte")
	}
	return rc.conn.Set(&item)
}

// Delete deletes a value in memcache.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.Delete(key)
}

// Incr increases counter.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Increment(key, 1)
	return err
}

// Decr decreases counter.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	_, err := rc.conn.Decrement(key, 1)
	return err
}

// IsExist checks if a value exists in memcache.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return false, err
		}
	}
	_, err := rc.conn.Get(key)
	return err == nil, err
}

// ClearAll clears all cache in memcache.
func (rc *Cache) ClearAll(context.Context) error {
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return rc.conn.FlushAll()
}

// StartAndGC starts the memcache adapter.
// config: must be in the format {"conn":"connection info"}.
// If an error occurs during connecting, an error is returned
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	rc.conninfo = strings.Split(cf["conn"], ";")
	if rc.conn == nil {
		if err := rc.connectInit(); err != nil {
			return err
		}
	}
	return nil
}

// connect to memcache and keep the connection.
func (rc *Cache) connectInit() error {
	rc.conn = memcache.New(rc.conninfo...)
	return nil
}

func init() {
	cache.Register("memcache", NewMemCache)
}
