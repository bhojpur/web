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
)

func TestNewHttpServerWithCfg(t *testing.T) {
	BConfig.AppName = "Before"
	svr := NewHttpServerWithCfg(BConfig)
	svr.Cfg.AppName = "hello"
	assert.Equal(t, "hello", BConfig.AppName)
}

func TestServerCtrlGet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/user", nil)
	w := httptest.NewRecorder()

	CtrlGet("/user", ExampleController.Ping)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlGet can't run")
	}
}

func TestServerCtrlPost(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, "/user", nil)
	w := httptest.NewRecorder()

	CtrlPost("/user", ExampleController.Ping)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlPost can't run")
	}
}

func TestServerCtrlHead(t *testing.T) {
	r, _ := http.NewRequest(http.MethodHead, "/user", nil)
	w := httptest.NewRecorder()

	CtrlHead("/user", ExampleController.Ping)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlHead can't run")
	}
}

func TestServerCtrlPut(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPut, "/user", nil)
	w := httptest.NewRecorder()

	CtrlPut("/user", ExampleController.Ping)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlPut can't run")
	}
}

func TestServerCtrlPatch(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPatch, "/user", nil)
	w := httptest.NewRecorder()

	CtrlPatch("/user", ExampleController.Ping)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlPatch can't run")
	}
}

func TestServerCtrlDelete(t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, "/user", nil)
	w := httptest.NewRecorder()

	CtrlDelete("/user", ExampleController.Ping)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestServerCtrlDelete can't run")
	}
}

func TestServerCtrlAny(t *testing.T) {
	CtrlAny("/user", ExampleController.Ping)

	for method := range HTTPMETHOD {
		r, _ := http.NewRequest(method, "/user", nil)
		w := httptest.NewRecorder()
		BhojpurApp.Handlers.ServeHTTP(w, r)
		if w.Body.String() != exampleBody {
			t.Errorf("TestServerCtrlAny can't run")
		}
	}
}
