package memcache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/bradfitz/gomemcache/memcache"

	"github.com/bhojpur/web/pkg/client/cache"
)

func TestMemcacheCache(t *testing.T) {
	addr := os.Getenv("MEMCACHE_ADDR")
	if addr == "" {
		addr = "127.0.0.1:11211"
	}

	bm, err := cache.NewCache("memcache", fmt.Sprintf(`{"conn": "%s"}`, addr))
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put(context.Background(), "bhojpur", "1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "bhojpur"); !res {
		t.Error("check err")
	}

	time.Sleep(11 * time.Second)

	if res, _ := bm.IsExist(context.Background(), "bhojpur"); res {
		t.Error("check err")
	}
	if err = bm.Put(context.Background(), "bhojpur", "1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	val, _ := bm.Get(context.Background(), "bhojpur")
	if v, err := strconv.Atoi(string(val.([]byte))); err != nil || v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr(context.Background(), "bhojpur"); err != nil {
		t.Error("Incr Error", err)
	}

	val, _ = bm.Get(context.Background(), "bhojpur")
	if v, err := strconv.Atoi(string(val.([]byte))); err != nil || v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr(context.Background(), "bhojpur"); err != nil {
		t.Error("Decr Error", err)
	}

	val, _ = bm.Get(context.Background(), "bhojpur")
	if v, err := strconv.Atoi(string(val.([]byte))); err != nil || v != 1 {
		t.Error("get err")
	}
	bm.Delete(context.Background(), "bhojpur")
	if res, _ := bm.IsExist(context.Background(), "bhojpur"); res {
		t.Error("delete err")
	}

	// test string
	if err = bm.Put(context.Background(), "bhojpur", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "bhojpur"); !res {
		t.Error("check err")
	}

	val, _ = bm.Get(context.Background(), "bhojpur")
	if v := val.([]byte); string(v) != "author" {
		t.Error("get err")
	}

	// test GetMulti
	if err = bm.Put(context.Background(), "bhojpur1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "bhojpur1"); !res {
		t.Error("check err")
	}

	vv, _ := bm.GetMulti(context.Background(), []string{"bhojpur", "bhojpur1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if string(vv[0].([]byte)) != "author" && string(vv[0].([]byte)) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if string(vv[1].([]byte)) != "author1" && string(vv[1].([]byte)) != "author" {
		t.Error("GetMulti ERROR")
	}

	vv, err = bm.GetMulti(context.Background(), []string{"bhojpur0", "bhojpur1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0] != nil {
		t.Error("GetMulti ERROR")
	}
	if string(vv[1].([]byte)) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if err != nil && err.Error() == "key [bhojpur0] error: key isn't exist" {
		t.Error("GetMulti ERROR")
	}

	// test clear all
	if err = bm.ClearAll(context.Background()); err != nil {
		t.Error("clear all err")
	}
}
