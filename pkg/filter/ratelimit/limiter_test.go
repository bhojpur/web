package ratelimit

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
	"time"

	"github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

func testRequest(t *testing.T, handler *websvr.ControllerRegister, requestIP, method, path string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	r.Header.Set("X-Real-Ip", requestIP)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != code {
		t.Errorf("%s, %s, %s: %d, supposed to be %d", requestIP, method, path, w.Code, code)
	}
}

func TestLimiter(t *testing.T) {
	handler := websvr.NewControllerRegister()
	err := handler.InsertFilter("/foo/*", websvr.BeforeRouter, NewLimiter(WithRate(1*time.Millisecond), WithCapacity(1), WithSessionKey(RemoteIPSessionKey)))
	if err != nil {
		t.Error(err)
	}
	handler.Any("*", func(ctx *context.Context) {
		ctx.Output.SetStatus(200)
	})

	route := "/foo/1"
	ip := "127.0.0.1"
	testRequest(t, handler, ip, "GET", route, 200)
	testRequest(t, handler, ip, "GET", route, 429)
	testRequest(t, handler, "127.0.0.2", "GET", route, 200)
	time.Sleep(1 * time.Millisecond)
	testRequest(t, handler, ip, "GET", route, 200)
}

func BenchmarkWithoutLimiter(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	websvr.BConfig.RunMode = websvr.PROD
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			handler.ServeHTTP(recorder, r)
		}
	})
}

func BenchmarkWithLimiter(b *testing.B) {
	recorder := httptest.NewRecorder()
	handler := websvr.NewControllerRegister()
	websvr.BConfig.RunMode = websvr.PROD
	err := handler.InsertFilter("*", websvr.BeforeRouter, NewLimiter(WithRate(1*time.Millisecond), WithCapacity(100)))
	if err != nil {
		b.Error(err)
	}
	handler.Any("/foo", func(ctx *context.Context) {
		ctx.Output.SetStatus(500)
	})
	b.ResetTimer()
	r, _ := http.NewRequest("PUT", "/foo", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			handler.ServeHTTP(recorder, r)
		}
	})
}
