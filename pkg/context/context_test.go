package context

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

	session "github.com/bhojpur/session/pkg/engine"
)

func TestXsrfReset_01(t *testing.T) {
	r := &http.Request{}
	c := NewContext()
	c.Request = r
	c.ResponseWriter = &Response{}
	c.ResponseWriter.reset(httptest.NewRecorder())
	c.Output.Reset(c)
	c.Input.Reset(c)
	c.XSRFToken("key", 16)
	if c._xsrfToken == "" {
		t.FailNow()
	}
	token := c._xsrfToken
	c.Reset(&Response{ResponseWriter: httptest.NewRecorder()}, r)
	if c._xsrfToken != "" {
		t.FailNow()
	}
	c.XSRFToken("key", 16)
	if c._xsrfToken == "" {
		t.FailNow()
	}
	if token == c._xsrfToken {
		t.FailNow()
	}
}

func TestContext_Session(t *testing.T) {
	c := NewContext()
	if store, err := c.Session(); store != nil || err == nil {
		t.FailNow()
	}
}

func TestContext_Session1(t *testing.T) {
	c := Context{}
	if store, err := c.Session(); store != nil || err == nil {
		t.FailNow()
	}
}

func TestContext_Session2(t *testing.T) {
	c := NewContext()
	c.Input.CruSession = &session.MemSessionStore{}

	if store, err := c.Session(); store == nil || err != nil {
		t.FailNow()
	}
}

func TestSetCookie(t *testing.T) {
	type cookie struct {
		Name     string
		Value    string
		MaxAge   int64
		Path     string
		Domain   string
		Secure   bool
		HttpOnly bool
		SameSite string
	}
	type testItem struct {
		item cookie
		want string
	}
	cases := []struct {
		request string
		valueGp []testItem
	}{
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "Strict"}, "name=value; Max-Age=0; Path=/; SameSite=Strict"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "Lax"}, "name=value; Max-Age=0; Path=/; SameSite=Lax"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, "None"}, "name=value; Max-Age=0; Path=/; SameSite=None"}}},
		{"/", []testItem{{cookie{"name", "value", -1, "/", "", false, false, ""}, "name=value; Max-Age=0; Path=/"}}},
	}
	for _, c := range cases {
		r, _ := http.NewRequest("GET", c.request, nil)
		output := NewOutput()
		output.Context = NewContext()
		output.Context.Reset(httptest.NewRecorder(), r)
		for _, item := range c.valueGp {
			params := item.item
			var others = []interface{}{params.MaxAge, params.Path, params.Domain, params.Secure, params.HttpOnly, params.SameSite}
			output.Context.SetCookie(params.Name, params.Value, others...)
			got := output.Context.ResponseWriter.Header().Get("Set-Cookie")
			if got != item.want {
				t.Fatalf("SetCookie error,should be:\n%v \ngot:\n%v", item.want, got)
			}
		}
	}
}
