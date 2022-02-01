package redis_sentinel

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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bhojpur/web/pkg/adapter/session"
)

func TestRedisSentinel(t *testing.T) {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "bsessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "127.0.0.1:6379,100,,0,master",
	}
	globalSessions, e := session.NewManager("redis_sentinel", sessionConfig)
	if e != nil {
		t.Log(e)
		return
	}
	// todo test if e==nil
	go globalSessions.GC()

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		t.Fatal("session start failed:", err)
	}
	defer sess.SessionRelease(w)

	// SET AND GET
	err = sess.Set("username", "bhojpur")
	if err != nil {
		t.Fatal("set username failed:", err)
	}
	username := sess.Get("username")
	if username != "bhojpur" {
		t.Fatal("get username failed")
	}

	// DELETE
	err = sess.Delete("username")
	if err != nil {
		t.Fatal("delete username failed:", err)
	}
	username = sess.Get("username")
	if username != nil {
		t.Fatal("delete username failed")
	}

	// FLUSH
	err = sess.Set("username", "bhojpur")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	err = sess.Set("password", "1qaz2wsx")
	if err != nil {
		t.Fatal("set failed:", err)
	}
	username = sess.Get("username")
	if username != "bhojpur" {
		t.Fatal("get username failed")
	}
	password := sess.Get("password")
	if password != "1qaz2wsx" {
		t.Fatal("get password failed")
	}
	err = sess.Flush()
	if err != nil {
		t.Fatal("flush failed:", err)
	}
	username = sess.Get("username")
	if username != nil {
		t.Fatal("flush failed")
	}
	password = sess.Get("password")
	if password != nil {
		t.Fatal("flush failed")
	}

	sess.SessionRelease(w)

}
