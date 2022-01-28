package redis

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/bhojpur/web/pkg/adapter/cache"
)

func TestRedisCache(t *testing.T) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "127.0.0.1:6379"
	}

	bm, err := cache.NewCache("redis", fmt.Sprintf(`{"conn": "%s"}`, redisAddr))
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("bhojpur", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("bhojpur") {
		t.Error("check err")
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist("bhojpur") {
		t.Error("check err")
	}
	if err = bm.Put("bhojpur", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if v, _ := redis.Int(bm.Get("bhojpur"), err); v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("bhojpur"); err != nil {
		t.Error("Incr Error", err)
	}

	if v, _ := redis.Int(bm.Get("bhojpur"), err); v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("bhojpur"); err != nil {
		t.Error("Decr Error", err)
	}

	if v, _ := redis.Int(bm.Get("bhojpur"), err); v != 1 {
		t.Error("get err")
	}
	bm.Delete("bhojpur")
	if bm.IsExist("bhojpur") {
		t.Error("delete err")
	}

	//test string
	if err = bm.Put("bhojpur", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("bhojpur") {
		t.Error("check err")
	}

	if v, _ := redis.String(bm.Get("bhojpur"), err); v != "author" {
		t.Error("get err")
	}

	//test GetMulti
	if err = bm.Put("bhojpur1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("bhojpur1") {
		t.Error("check err")
	}

	vv := bm.GetMulti([]string{"bhojpur", "bhojpur1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if v, _ := redis.String(vv[0], nil); v != "author" {
		t.Error("GetMulti ERROR")
	}
	if v, _ := redis.String(vv[1], nil); v != "author1" {
		t.Error("GetMulti ERROR")
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}

func TestCache_Scan(t *testing.T) {
	timeoutDuration := 10 * time.Second
	// init
	bm, err := cache.NewCache("redis", `{"conn": "127.0.0.1:6379"}`)
	if err != nil {
		t.Error("init err")
	}
	// insert all
	for i := 0; i < 10000; i++ {
		if err = bm.Put(fmt.Sprintf("bhojpur%d", i), fmt.Sprintf("author%d", i), timeoutDuration); err != nil {
			t.Error("set Error", err)
		}
	}

	// clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}

}
