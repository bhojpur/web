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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	ctxsvr "github.com/bhojpur/web/pkg/context"
)

func TestControllerRegister_InsertFilterChain(t *testing.T) {

	InsertFilterChain("/*", func(next FilterFunc) FilterFunc {
		return func(ctx *ctxsvr.Context) {
			ctx.Output.Header("filter", "filter-chain")
			next(ctx)
		}
	})

	ns := NewNamespace("/chain")

	ns.Get("/*", func(ctx *ctxsvr.Context) {
		ctx.Output.Body([]byte("hello"))
	})

	r, _ := http.NewRequest("GET", "/chain/user", nil)
	w := httptest.NewRecorder()

	BhojpurApp.Handlers.ServeHTTP(w, r)

	assert.Equal(t, "filter-chain", w.Header().Get("filter"))
}
