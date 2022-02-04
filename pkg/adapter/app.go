package adapter

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

	context2 "github.com/bhojpur/web/pkg/adapter/context"
	"github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

// BhojpurApp is an application instance
var BhojpurApp *App

func init() {
	// create Bhojpur.NET Platform application
	BhojpurApp = (*App)(websvr.BhojpurApp)
}

// App defines Bhojpur.NET Platform application with a new PatternServeMux.
type App websvr.HttpServer

// NewApp returns a new Bhojpur.NET Platform application.
func NewApp() *App {
	return (*App)(websvr.NewHttpSever())
}

// MiddleWare function for http.Handler
type MiddleWare websvr.MiddleWare

// Run Bhojpur.NET Platform application.
func (app *App) Run(mws ...MiddleWare) {
	newMws := oldMiddlewareToNew(mws)
	(*websvr.HttpServer)(app).Run("", newMws...)
}

func oldMiddlewareToNew(mws []MiddleWare) []websvr.MiddleWare {
	newMws := make([]websvr.MiddleWare, 0, len(mws))
	for _, old := range mws {
		newMws = append(newMws, (websvr.MiddleWare)(old))
	}
	return newMws
}

// Router adds a patterned controller handler to BhojpurApp.
// it's an alias method of HttpServer.Router.
// usage:
//  simple router
//  websvr.Router("/admin", &admin.UserController{})
//  websvr.Router("/admin/index", &admin.ArticleController{})
//
//  regex router
//
//  websvr.Router("/api/:id([0-9]+)", &controllers.RController{})
//
//  custom rules
//  websvr.Router("/api/list",&RestController{},"*:ListFood")
//  websvr.Router("/api/create",&RestController{},"post:CreateFood")
//  websvr.Router("/api/update",&RestController{},"put:UpdateFood")
//  websvr.Router("/api/delete",&RestController{},"delete:DeleteFood")
func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *App {
	return (*App)(websvr.Router(rootpath, c, mappingMethods...))
}

// UnregisterFixedRoute unregisters the route with the specified fixedRoute. It is particularly useful
// in web applications that inherit most routes from a base webapp via the underscore
// import, and aim to overwrite only certain paths.
// The method parameter can be empty or "*" for all HTTP methods, or a particular
// method type (e.g. "GET" or "POST") for selective removal.
//
// Usage (replace "GET" with "*" for all methods):
//  websvr.UnregisterFixedRoute("/yourpreviouspath", "GET")
//  websvr.Router("/yourpreviouspath", yourControllerAddress, "get:GetNewPage")
func UnregisterFixedRoute(fixedRoute string, method string) *App {
	return (*App)(websvr.UnregisterFixedRoute(fixedRoute, method))
}

// Include will generate router file in the router/xxx.go from the controller's comments
// usage:
// websvr.Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
// type BankAccount struct{
//   websvr.Controller
// }
//
// register the function
// func (b *BankAccount)Mapping(){
//  b.Mapping("ShowAccount" , b.ShowAccount)
//  b.Mapping("ModifyAccount", b.ModifyAccount)
// }
//
// //@router /account/:id  [get]
// func (b *BankAccount) ShowAccount(){
//    //logic
// }
//
//
// //@router /account/:id  [post]
// func (b *BankAccount) ModifyAccount(){
//    //logic
// }
//
// the comments @router url methodlist
// url support all the function Router's pattern
// methodlist [get post head put delete options *]
func Include(cList ...ControllerInterface) *App {
	newList := oldToNewCtrlIntfs(cList)
	return (*App)(websvr.Include(newList...))
}

func oldToNewCtrlIntfs(cList []ControllerInterface) []websvr.ControllerInterface {
	newList := make([]websvr.ControllerInterface, 0, len(cList))
	for _, c := range cList {
		newList = append(newList, c)
	}
	return newList
}

// RESTRouter adds a restful controller handler to BhojpurApp.
// its' controller implements websvr.ControllerInterface and
// defines a param "pattern/:objectId" to visit each resource.
func RESTRouter(rootpath string, c ControllerInterface) *App {
	return (*App)(websvr.RESTRouter(rootpath, c))
}

// AutoRouter adds defined controller handler to BhojpurApp.
// it's same to HttpServer.AutoRouter.
// if websvr.AddAuto(&MainController{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func AutoRouter(c ControllerInterface) *App {
	return (*App)(websvr.AutoRouter(c))
}

// AutoPrefix adds controller handler to BhojpurApp with prefix.
// it's same to HttpServer.AutoRouterWithPrefix.
// if websvr.AutoPrefix("/admin",&MainController{}) and MainController has methods List and Page,
// visit the url /admin/main/list to exec List function or /admin/main/page to exec Page function.
func AutoPrefix(prefix string, c ControllerInterface) *App {
	return (*App)(websvr.AutoPrefix(prefix, c))
}

// Get used to register router for Get method
// usage:
//    websvr.Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Get(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Get(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Post used to register router for Post method
// usage:
//    websvr.Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Post(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Post(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Delete used to register router for Delete method
// usage:
//    websvr.Delete("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Delete(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Delete(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Put used to register router for Put method
// usage:
//    websvr.Put("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Put(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Put(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Head used to register router for Head method
// usage:
//    websvr.Head("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Head(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Head(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Options used to register router for Options method
// usage:
//    websvr.Options("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Options(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Options(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Patch used to register router for Patch method
// usage:
//    websvr.Patch("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Patch(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Patch(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Any used to register router for all methods
// usage:
//    websvr.Any("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func Any(rootpath string, f FilterFunc) *App {
	return (*App)(websvr.Any(rootpath, func(ctx *context.Context) {
		f((*context2.Context)(ctx))
	}))
}

// Handler used to register a Handler router
// usage:
//    websvr.Handler("/api", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
//          fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//    }))
func Handler(rootpath string, h http.Handler, options ...interface{}) *App {
	return (*App)(websvr.Handler(rootpath, h, options))
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// websvr.BeforeStatic, websvr.BeforeRouter, websvr.BeforeExec, websvr.AfterExec and websvr.FinishRouter.
// The bool params is for setting the returnOnOutput value (false allows multiple filters to execute)
func InsertFilter(pattern string, pos int, filter FilterFunc, params ...bool) *App {
	opts := oldToNewFilterOpts(params)
	return (*App)(websvr.InsertFilter(pattern, pos, func(ctx *context.Context) {
		filter((*context2.Context)(ctx))
	}, opts...))
}
