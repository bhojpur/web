package redis

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

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/bhojpur/web/pkg/client/cache"
)

var (
	// The collection name of redis for the cache adapter.
	DefaultKey = "bcacheRedis"
)

// Cache is Redis cache adapter.
type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbNum    int
	key      string
	password string
	maxIdle  int

	// Timeout value (less than the redis server's timeout value)
	timeout time.Duration
}

// NewRedisCache creates a new redis cache with default collection name.
func NewRedisCache() cache.Cache {
	return &Cache{key: DefaultKey}
}

// Execute the redis commands. args[0] must be the key name
func (rc *Cache) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if len(args) < 1 {
		return nil, errors.New("missing required arguments")
	}
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	if v, err := rc.do("GET", key); err == nil {
		return v, nil
	} else {
		return nil, err
	}
}

// GetMulti gets cache from redis.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	c := rc.p.Get()
	defer c.Close()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	return redis.Values(c.Do("MGET", args...))
}

// Put puts cache into redis.
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete deletes a key's cache in redis.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist checks cache's existence in redis.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return v, nil
}

// Incr increases a key's counter in redis.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decreases a key's counter in redis.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// ClearAll deletes all cache in the redis collection
func (rc *Cache) ClearAll(context.Context) error {
	cachedKeys, err := rc.Scan(rc.key + ":*")
	if err != nil {
		return err
	}
	c := rc.p.Get()
	defer c.Close()
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// Scan scans all keys matching a given pattern.
func (rc *Cache) Scan(pattern string) (keys []string, err error) {
	c := rc.p.Get()
	defer c.Close()
	var (
		cursor uint64 = 0 // start
		result []interface{}
		list   []string
	)
	for {
		result, err = redis.Values(c.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 1024))
		if err != nil {
			return
		}
		list, err = redis.Strings(result[1], nil)
		if err != nil {
			return
		}
		keys = append(keys, list...)
		cursor, err = redis.Uint64(result[0], nil)
		if err != nil {
			return
		}
		if cursor == 0 { // over
			return
		}
	}
}

// StartAndGC starts the redis cache adapter.
// config: must be in this format {"key":"collection key","conn":"connection info","dbNum":"0"}
// Cached items in redis are stored forever, no garbage collection happens
func (rc *Cache) StartAndGC(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)

	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}

	// Format redis://<password>@<host>:<port>
	cf["conn"] = strings.Replace(cf["conn"], "redis://", "", 1)
	if i := strings.Index(cf["conn"], "@"); i > -1 {
		cf["password"] = cf["conn"][0:i]
		cf["conn"] = cf["conn"][i+1:]
	}

	if _, ok := cf["dbNum"]; !ok {
		cf["dbNum"] = "0"
	}
	if _, ok := cf["password"]; !ok {
		cf["password"] = ""
	}
	if _, ok := cf["maxIdle"]; !ok {
		cf["maxIdle"] = "3"
	}
	if _, ok := cf["timeout"]; !ok {
		cf["timeout"] = "180s"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbNum, _ = strconv.Atoi(cf["dbNum"])
	rc.password = cf["password"]
	rc.maxIdle, _ = strconv.Atoi(cf["maxIdle"])

	if v, err := time.ParseDuration(cf["timeout"]); err == nil {
		rc.timeout = v
	} else {
		rc.timeout = 180 * time.Second
	}

	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *Cache) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		if err != nil {
			return nil, err
		}

		if rc.password != "" {
			if _, err := c.Do("AUTH", rc.password); err != nil {
				c.Close()
				return nil, err
			}
		}

		_, selecterr := c.Do("SELECT", rc.dbNum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     rc.maxIdle,
		IdleTimeout: rc.timeout,
		Dial:        dialFunc,
	}
}

func init() {
	cache.Register("redis", NewRedisCache)
}
