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

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/bhojpur/web/pkg/adapter/cache"
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
	if err = bm.Put("bhojpur", "1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("bhojpur") {
		t.Error("check err")
	}

	time.Sleep(11 * time.Second)

	if bm.IsExist("bhojpur") {
		t.Error("check err")
	}
	if err = bm.Put("bhojpur", "1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get("bhojpur").([]byte))); err != nil || v != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("bhojpur"); err != nil {
		t.Error("Incr Error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get("bhojpur").([]byte))); err != nil || v != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("bhojpur"); err != nil {
		t.Error("Decr Error", err)
	}

	if v, err := strconv.Atoi(string(bm.Get("bhojpur").([]byte))); err != nil || v != 1 {
		t.Error("get err")
	}
	bm.Delete("bhojpur")
	if bm.IsExist("bhojpur") {
		t.Error("delete err")
	}

	// test string
	if err = bm.Put("bhojpur", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("bhojpur") {
		t.Error("check err")
	}

	if v := bm.Get("bhojpur").([]byte); string(v) != "author" {
		t.Error("get err")
	}

	// test GetMulti
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
	if string(vv[0].([]byte)) != "author" && string(vv[0].([]byte)) != "author1" {
		t.Error("GetMulti ERROR")
	}
	if string(vv[1].([]byte)) != "author1" && string(vv[1].([]byte)) != "author" {
		t.Error("GetMulti ERROR")
	}

	// test clear all
	if err = bm.ClearAll(); err != nil {
		t.Error("clear all err")
	}
}
