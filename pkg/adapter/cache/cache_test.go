package cache

import (
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
	//timeoutDuration := 10 * time.Second

	bm.Put("edwardhey", 0, time.Second*20)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			bm.Incr("edwardhey")
		}()
	}
	wg.Wait()
	if bm.Get("edwardhey").(int) != 10 {
		t.Error("Incr err")
	}
}

func TestCache(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
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

	if v := bm.Get("bhojpur"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("bhojpur") {
		t.Error("check err")
	}

	if err = bm.Put("bhojpur", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("bhojpur"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("bhojpur"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("bhojpur"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("bhojpur"); v.(int) != 1 {
		t.Error("get err")
	}
	bm.Delete("bhojpur")
	if bm.IsExist("bhojpur") {
		t.Error("delete err")
	}

	//test GetMulti
	if err = bm.Put("bhojpur", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("bhojpur") {
		t.Error("check err")
	}
	if v := bm.Get("bhojpur"); v.(string) != "author" {
		t.Error("get err")
	}

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
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache("file", `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)
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

	if v := bm.Get("bhojpur"); v.(int) != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("bhojpur"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("bhojpur"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("bhojpur"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("bhojpur"); v.(int) != 1 {
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
	if v := bm.Get("bhojpur"); v.(string) != "author" {
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
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}

	os.RemoveAll("cache")
}
