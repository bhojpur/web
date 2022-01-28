package cache

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"
)

func TestCacheIncr(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	// timeoutDuration := 10 * time.Second

	bm.Put(context.Background(), "edwardhey", 0, time.Second*20)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			bm.Incr(context.Background(), "edwardhey")
		}()
	}
	wg.Wait()
	val, _ := bm.Get(context.Background(), "edwardhey")
	if val.(int) != 10 {
		t.Error("Incr err")
	}
}

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put(context.Background(), "bhojpur", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "bhojpur"); !res {
		t.Error("check err")
	}

	if v, _ := bm.Get(context.Background(), "bhojpur"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if res, _ := bm.IsExist(context.Background(), "bhojpur"); res {
		t.Error("check err")
	}

	if err = bm.Put(context.Background(), "bhojpur", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	// test different integer type for incr & decr
	testMultiIncrDecr(t, bm, timeoutDuration)

	bm.Delete(context.Background(), "bhojpur")
	if res, _ := bm.IsExist(context.Background(), "bhojpur"); res {
		t.Error("delete err")
	}

	// test GetMulti
	if err = bm.Put(context.Background(), "bhojpur", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "bhojpur"); !res {
		t.Error("check err")
	}
	if v, _ := bm.Get(context.Background(), "bhojpur"); v.(string) != "author" {
		t.Error("get err")
	}

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
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}

	vv, err = bm.GetMulti(context.Background(), []string{"bhojpur0", "bhojpur1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0] != nil {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if err != nil && err.Error() != "key [bhojpur0] error: the key isn't exist" {
		t.Error("GetMulti ERROR")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put(context.Background(), "bhojpur", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if res, _ := bm.IsExist(context.Background(), "bhojpur"); !res {
		t.Error("check err")
	}

	if v, _ := bm.Get(context.Background(), "bhojpur"); v.(int) != 1 {
		t.Error("get err")
	}

	// test different integer type for incr & decr
	testMultiIncrDecr(t, bm, timeoutDuration)

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
	if v, _ := bm.Get(context.Background(), "bhojpur"); v.(string) != "author" {
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
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}

	vv, err = bm.GetMulti(context.Background(), []string{"bhojpur0", "bhojpur1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0] != nil {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if err == nil {
		t.Error("GetMulti ERROR")
	}

	os.RemoveAll("cache")
}

func testMultiIncrDecr(t *testing.T, c Cache, timeout time.Duration) {
	testIncrDecr(t, c, 1, 2, timeout)
	testIncrDecr(t, c, int32(1), int32(2), timeout)
	testIncrDecr(t, c, int64(1), int64(2), timeout)
	testIncrDecr(t, c, uint(1), uint(2), timeout)
	testIncrDecr(t, c, uint32(1), uint32(2), timeout)
	testIncrDecr(t, c, uint64(1), uint64(2), timeout)
}

func testIncrDecr(t *testing.T, c Cache, beforeIncr interface{}, afterIncr interface{}, timeout time.Duration) {
	var err error
	ctx := context.Background()
	key := "incDecKey"
	if err = c.Put(ctx, key, beforeIncr, timeout); err != nil {
		t.Error("Get Error", err)
	}

	if err = c.Incr(ctx, key); err != nil {
		t.Error("Incr Error", err)
	}

	if v, _ := c.Get(ctx, key); v != afterIncr {
		t.Error("Get Error")
	}

	if err = c.Decr(ctx, key); err != nil {
		t.Error("Decr Error", err)
	}

	if v, _ := c.Get(ctx, key); v != beforeIncr {
		t.Error("Get Error")
	}

	if err := c.Delete(ctx, key); err != nil {
		t.Error("Delete Error")
	}
}
