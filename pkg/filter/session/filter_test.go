package session

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

	"github.com/google/uuid"

	session "github.com/bhojpur/session/pkg/engine"
	webContext "github.com/bhojpur/web/pkg/context"
	web "github.com/bhojpur/web/pkg/engine"
)

func testRequest(t *testing.T, handler *web.ControllerRegister, path string, method string, code int) {
	r, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != code {
		t.Errorf("%s, %s: %d, supposed to be %d", path, method, w.Code, code)
	}
}

func TestSession(t *testing.T) {
	storeKey := uuid.New().String()
	handler := web.NewControllerRegister()
	handler.InsertFilterChain(
		"*",
		Session(
			session.ProviderMemory,
			session.CfgCookieName(`go_session_id`),
			session.CfgSetCookie(true),
			session.CfgGcLifeTime(3600),
			session.CfgMaxLifeTime(3600),
			session.CfgSecure(false),
			session.CfgCookieLifeTime(3600),
		),
	)
	handler.InsertFilterChain(
		"*",
		func(next web.FilterFunc) web.FilterFunc {
			return func(ctx *webContext.Context) {
				if store := ctx.Input.GetData(storeKey); store == nil {
					t.Error(`store should not be nil`)
				}
				next(ctx)
			}
		},
	)
	handler.Any("*", func(ctx *webContext.Context) {
		ctx.Output.SetStatus(200)
	})

	testRequest(t, handler, "/dataset1/resource1", "GET", 200)
}

func TestSession1(t *testing.T) {
	handler := web.NewControllerRegister()
	handler.InsertFilterChain(
		"*",
		Session(
			session.ProviderMemory,
			session.CfgCookieName(`go_session_id`),
			session.CfgSetCookie(true),
			session.CfgGcLifeTime(3600),
			session.CfgMaxLifeTime(3600),
			session.CfgSecure(false),
			session.CfgCookieLifeTime(3600),
		),
	)
	handler.InsertFilterChain(
		"*",
		func(next web.FilterFunc) web.FilterFunc {
			return func(ctx *webContext.Context) {
				if store, err := ctx.Session(); store == nil || err != nil {
					t.Error(`store should not be nil`)
				}
				next(ctx)
			}
		},
	)
	handler.Any("*", func(ctx *webContext.Context) {
		ctx.Output.SetStatus(200)
	})

	testRequest(t, handler, "/dataset1/resource1", "GET", 200)
}
