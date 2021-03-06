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

	"github.com/bhojpur/web/pkg/context"
)

const (
	exampleBody        = "hello world"
	examplePointerBody = "hello world pointer"

	nsNamespace     = "/router"
	nsPath          = "/user"
	nsNamespacePath = "/router/user"
)

type ExampleController struct {
	Controller
}

func (m ExampleController) Ping() {
	err := m.Ctx.Output.Body([]byte(exampleBody))
	if err != nil {
		fmt.Println(err)
	}
}

func (m *ExampleController) PingPointer() {
	err := m.Ctx.Output.Body([]byte(examplePointerBody))
	if err != nil {
		fmt.Println(err)
	}
}

func (m ExampleController) ping() {
	err := m.Ctx.Output.Body([]byte("ping method"))
	if err != nil {
		fmt.Println(err)
	}
}

func TestNamespaceGet(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/user", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Get("/user", func(ctx *context.Context) {
		ctx.Output.Body([]byte("v1_user"))
	})
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "v1_user" {
		t.Errorf("TestNamespaceGet can't run, get the response is " + w.Body.String())
	}
}

func TestNamespacePost(t *testing.T) {
	r, _ := http.NewRequest("POST", "/v1/user/123", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Post("/user/:id", func(ctx *context.Context) {
		ctx.Output.Body([]byte(ctx.Input.Param(":id")))
	})
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "123" {
		t.Errorf("TestNamespacePost can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNest(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/admin/order", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Namespace(
		NewNamespace("/admin").
			Get("/order", func(ctx *context.Context) {
				ctx.Output.Body([]byte("order"))
			}),
	)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "order" {
		t.Errorf("TestNamespaceNest can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNestParam(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/admin/order/123", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Namespace(
		NewNamespace("/admin").
			Get("/order/:id", func(ctx *context.Context) {
				ctx.Output.Body([]byte(ctx.Input.Param(":id")))
			}),
	)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "123" {
		t.Errorf("TestNamespaceNestParam can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceRouter(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/api/list", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Router("/api/list", &TestController{}, "*:List")
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("TestNamespaceRouter can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceAutoFunc(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/test/list", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.AutoRouter(&TestController{})
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "i am list" {
		t.Errorf("user define func can't run")
	}
}

func TestNamespaceFilter(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v1/user/123", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v1")
	ns.Filter("before", func(ctx *context.Context) {
		ctx.Output.Body([]byte("this is Filter"))
	}).
		Get("/user/:id", func(ctx *context.Context) {
			ctx.Output.Body([]byte(ctx.Input.Param(":id")))
		})
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "this is Filter" {
		t.Errorf("TestNamespaceFilter can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCond(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v2/test/list", nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/v2")
	ns.Cond(func(ctx *context.Context) bool {
		return ctx.Input.Domain() == "bhojpur.net"
	}).
		AutoRouter(&TestController{})
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Code != 405 {
		t.Errorf("TestNamespaceCond can't run get the result " + strconv.Itoa(w.Code))
	}
}

func TestNamespaceInside(t *testing.T) {
	r, _ := http.NewRequest("GET", "/v3/shop/order/123", nil)
	w := httptest.NewRecorder()
	ns := NewNamespace("/v3",
		NSAutoRouter(&TestController{}),
		NSNamespace("/shop",
			NSGet("/order/:id", func(ctx *context.Context) {
				ctx.Output.Body([]byte(ctx.Input.Param(":id")))
			}),
		),
	)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != "123" {
		t.Errorf("TestNamespaceInside can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlGet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	ns.CtrlGet(nsPath, ExampleController.Ping)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceCtrlGet can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlPost(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	ns.CtrlPost(nsPath, ExampleController.Ping)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceCtrlPost can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlDelete(t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	ns.CtrlDelete(nsPath, ExampleController.Ping)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceCtrlDelete can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlPut(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPut, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	ns.CtrlPut(nsPath, ExampleController.Ping)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceCtrlPut can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlHead(t *testing.T) {
	r, _ := http.NewRequest(http.MethodHead, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	ns.CtrlHead(nsPath, ExampleController.Ping)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceCtrlHead can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlOptions(t *testing.T) {
	r, _ := http.NewRequest(http.MethodOptions, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	ns.CtrlOptions(nsPath, ExampleController.Ping)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceCtrlOptions can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlPatch(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPatch, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	ns.CtrlPatch(nsPath, ExampleController.Ping)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceCtrlPatch can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceCtrlAny(t *testing.T) {
	ns := NewNamespace(nsNamespace)
	ns.CtrlAny(nsPath, ExampleController.Ping)
	AddNamespace(ns)

	for method := range HTTPMETHOD {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(method, nsNamespacePath, nil)
		BhojpurApp.Handlers.ServeHTTP(w, r)
		if w.Body.String() != exampleBody {
			t.Errorf("TestNamespaceCtrlAny can't run, get the response is " + w.Body.String())
		}
	}
}

func TestNamespaceNSCtrlGet(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	NSCtrlGet(nsPath, ExampleController.Ping)(ns)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceNSCtrlGet can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSCtrlPost(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPost, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace("/router")
	NSCtrlPost(nsPath, ExampleController.Ping)(ns)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceNSCtrlPost can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSCtrlDelete(t *testing.T) {
	r, _ := http.NewRequest(http.MethodDelete, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	NSCtrlDelete(nsPath, ExampleController.Ping)(ns)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceNSCtrlDelete can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSCtrlPut(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPut, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	NSCtrlPut(nsPath, ExampleController.Ping)(ns)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceNSCtrlPut can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSCtrlHead(t *testing.T) {
	r, _ := http.NewRequest(http.MethodHead, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	NSCtrlHead(nsPath, ExampleController.Ping)(ns)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceNSCtrlHead can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSCtrlOptions(t *testing.T) {
	r, _ := http.NewRequest(http.MethodOptions, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	NSCtrlOptions(nsPath, ExampleController.Ping)(ns)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceNSCtrlOptions can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSCtrlPatch(t *testing.T) {
	r, _ := http.NewRequest(http.MethodPatch, nsNamespacePath, nil)
	w := httptest.NewRecorder()

	ns := NewNamespace(nsNamespace)
	NSCtrlPatch("/user", ExampleController.Ping)(ns)
	AddNamespace(ns)
	BhojpurApp.Handlers.ServeHTTP(w, r)
	if w.Body.String() != exampleBody {
		t.Errorf("TestNamespaceNSCtrlPatch can't run, get the response is " + w.Body.String())
	}
}

func TestNamespaceNSCtrlAny(t *testing.T) {
	ns := NewNamespace(nsNamespace)
	NSCtrlAny(nsPath, ExampleController.Ping)(ns)
	AddNamespace(ns)

	for method := range HTTPMETHOD {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest(method, nsNamespacePath, nil)
		BhojpurApp.Handlers.ServeHTTP(w, r)
		if w.Body.String() != exampleBody {
			t.Errorf("TestNamespaceNSCtrlAny can't run, get the response is " + w.Body.String())
		}
	}
}
