package engine

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/bhojpur/web/pkg/context"
)

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
		return ctx.Input.Domain() == "beego.me"
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
