package engine

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
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/pkg/context"
)

func TestControllerRegisterInsertFilterChain(t *testing.T) {
	InsertFilterChain("/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", "filter-chain")
			next(ctx)
		}
	})

	ns := NewNamespace("/chain")

	ns.Get("/*", func(ctx *context.Context) {
		_ = ctx.Output.Body([]byte("hello"))
	})

	r, _ := http.NewRequest("GET", "/chain/user", nil)
	w := httptest.NewRecorder()

	BhojpurApp.Handlers.Init()
	BhojpurApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, "filter-chain", w.Header().Get("filter"))
}

func TestControllerRegister_InsertFilterChain_Order(t *testing.T) {
	InsertFilterChain("/abc", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("first", fmt.Sprintf("%d", time.Now().UnixNano()))
			time.Sleep(time.Millisecond * 10)
			next(ctx)
		}
	})

	InsertFilterChain("/abc", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("second", fmt.Sprintf("%d", time.Now().UnixNano()))
			time.Sleep(time.Millisecond * 10)
			next(ctx)
		}
	})

	r, _ := http.NewRequest("GET", "/abc", nil)
	w := httptest.NewRecorder()

	BhojpurApp.Handlers.Init()
	BhojpurApp.Handlers.ServeHTTP(w, r)
	first := w.Header().Get("first")
	second := w.Header().Get("second")

	ft, _ := strconv.ParseInt(first, 10, 64)
	st, _ := strconv.ParseInt(second, 10, 64)

	assert.True(t, st > ft)
}

func TestFilterChainRouter(t *testing.T) {
	app := NewHttpSever()
	const filterNonMatch = "filter-chain-non-match"
	app.InsertFilterChain("/app/nonMatch/before/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", filterNonMatch)
			next(ctx)
		}
	})

	const filterAll = "filter-chain-all"
	app.InsertFilterChain("/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", filterAll)
			next(ctx)
		}
	})

	app.InsertFilterChain("/app/nonMatch/after/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("filter", filterNonMatch)
			next(ctx)
		}
	})

	app.InsertFilterChain("/app/match/*", func(next FilterFunc) FilterFunc {
		return func(ctx *context.Context) {
			ctx.Output.Header("match", "yes")
			next(ctx)
		}
	})

	app.Handlers.Init()

	r, _ := http.NewRequest("GET", "/app/match", nil)
	w := httptest.NewRecorder()

	app.Handlers.ServeHTTP(w, r)
	assert.Equal(t, filterAll, w.Header().Get("filter"))
	assert.Equal(t, "yes", w.Header().Get("match"))

	r, _ = http.NewRequest("GET", "/app/match1", nil)
	w = httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, r)
	assert.Equal(t, filterAll, w.Header().Get("filter"))
	assert.NotEqual(t, "yes", w.Header().Get("match"))

	r, _ = http.NewRequest("GET", "/app/nonMatch", nil)
	w = httptest.NewRecorder()
	app.Handlers.ServeHTTP(w, r)
	assert.Equal(t, filterAll, w.Header().Get("filter"))
	assert.NotEqual(t, "yes", w.Header().Get("match"))
}
