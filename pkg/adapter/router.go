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
	"time"

	ctxadp "github.com/bhojpur/web/pkg/adapter/context"
	"github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

// default filter execution points
const (
	BeforeStatic = websvr.BeforeStatic
	BeforeRouter = websvr.BeforeRouter
	BeforeExec   = websvr.BeforeExec
	AfterExec    = websvr.AfterExec
	FinishRouter = websvr.FinishRouter
)

var (
	// HTTPMETHOD list the supported http methods.
	HTTPMETHOD = websvr.HTTPMETHOD

	// DefaultAccessLogFilter will skip the accesslog if return true
	DefaultAccessLogFilter FilterHandler = &newToOldFtHdlAdapter{
		delegate: websvr.DefaultAccessLogFilter,
	}
)

// FilterHandler is an interface for
type FilterHandler interface {
	Filter(*ctxadp.Context) bool
}

type newToOldFtHdlAdapter struct {
	delegate websvr.FilterHandler
}

func (n *newToOldFtHdlAdapter) Filter(ctx *ctxadp.Context) bool {
	return n.delegate.Filter((*context.Context)(ctx))
}

// ExceptMethodAppend to append a slice's value into "exceptMethod", for controller's methods shouldn't reflect to AutoRouter
func ExceptMethodAppend(action string) {
	websvr.ExceptMethodAppend(action)
}

// ControllerInfo holds information about the controller.
type ControllerInfo websvr.ControllerInfo

func (c *ControllerInfo) GetPattern() string {
	return (*websvr.ControllerInfo)(c).GetPattern()
}

// ControllerRegister containers registered router rules, controller handlers and filters.
type ControllerRegister websvr.ControllerRegister

// NewControllerRegister returns a new ControllerRegister.
func NewControllerRegister() *ControllerRegister {
	return (*ControllerRegister)(websvr.NewControllerRegister())
}

// Add controller handler and pattern rules to ControllerRegister.
// usage:
//	default methods is the same name as method
//	Add("/user",&UserController{})
//	Add("/api/list",&RestController{},"*:ListFood")
//	Add("/api/create",&RestController{},"post:CreateFood")
//	Add("/api/update",&RestController{},"put:UpdateFood")
//	Add("/api/delete",&RestController{},"delete:DeleteFood")
//	Add("/api",&RestController{},"get,post:ApiFunc"
//	Add("/simple",&SimpleController{},"get:GetFunc;post:PostFunc")
func (p *ControllerRegister) Add(pattern string, c ControllerInterface, mappingMethods ...string) {
	(*websvr.ControllerRegister)(p).Add(pattern, c, websvr.WithRouterMethods(c, mappingMethods...))
}

// Include only when the Runmode is dev will generate router file in the router/auto.go from the controller
// Include(&BankAccount{}, &OrderController{},&RefundController{},&ReceiptController{})
func (p *ControllerRegister) Include(cList ...ControllerInterface) {
	nls := oldToNewCtrlIntfs(cList)
	(*websvr.ControllerRegister)(p).Include(nls...)
}

// GetContext returns a context from pool, so usually you should remember to call Reset function to clean the context
// And don't forget to give back context to pool
// example:
//  ctx := p.GetContext()
//  ctx.Reset(w, q)
//  defer p.GiveBackContext(ctx)
func (p *ControllerRegister) GetContext() *ctxadp.Context {
	return (*ctxadp.Context)((*websvr.ControllerRegister)(p).GetContext())
}

// GiveBackContext put the ctx into pool so that it could be reuse
func (p *ControllerRegister) GiveBackContext(ctx *ctxadp.Context) {
	(*websvr.ControllerRegister)(p).GiveBackContext((*context.Context)(ctx))
}

// Get add get method
// usage:
//    Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Get(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Get(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Post add post method
// usage:
//    Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Post(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Post(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Put add put method
// usage:
//    Put("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Put(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Put(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Delete add delete method
// usage:
//    Delete("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Delete(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Delete(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Head add head method
// usage:
//    Head("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Head(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Head(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Patch add patch method
// usage:
//    Patch("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Patch(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Patch(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Options add options method
// usage:
//    Options("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Options(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Options(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Any add all method
// usage:
//    Any("/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) Any(pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).Any(pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// AddMethod add http method router
// usage:
//    AddMethod("get","/api/:id", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (p *ControllerRegister) AddMethod(method, pattern string, f FilterFunc) {
	(*websvr.ControllerRegister)(p).AddMethod(method, pattern, func(ctx *context.Context) {
		f((*ctxadp.Context)(ctx))
	})
}

// Handler add user defined Handler
func (p *ControllerRegister) Handler(pattern string, h http.Handler, options ...interface{}) {
	(*websvr.ControllerRegister)(p).Handler(pattern, h, options)
}

// AddAuto router to ControllerRegister.
// example websvr.AddAuto(&MainController{}),
// MainController has method List and Page.
// visit the url /main/list to execute List function
// /main/page to execute Page function.
func (p *ControllerRegister) AddAuto(c ControllerInterface) {
	(*websvr.ControllerRegister)(p).AddAuto(c)
}

// AddAutoPrefix Add auto router to ControllerRegister with prefix.
// example websvr.AddAutoPrefix("/admin",&MainController{}),
// MainController has method List and Page.
// visit the url /admin/main/list to execute List function
// /admin/main/page to execute Page function.
func (p *ControllerRegister) AddAutoPrefix(prefix string, c ControllerInterface) {
	(*websvr.ControllerRegister)(p).AddAutoPrefix(prefix, c)
}

// InsertFilter Add a FilterFunc with pattern rule and action constant.
// params is for:
//   1. setting the returnOnOutput value (false allows multiple filters to execute)
//   2. determining whether or not params need to be reset.
func (p *ControllerRegister) InsertFilter(pattern string, pos int, filter FilterFunc, params ...bool) error {
	opts := oldToNewFilterOpts(params)
	return (*websvr.ControllerRegister)(p).InsertFilter(pattern, pos, func(ctx *context.Context) {
		filter((*ctxadp.Context)(ctx))
	}, opts...)
}

func oldToNewFilterOpts(params []bool) []websvr.FilterOpt {
	opts := make([]websvr.FilterOpt, 0, 4)
	if len(params) > 0 {
		opts = append(opts, websvr.WithReturnOnOutput(params[0]))
	} else {
		// the default value should be true
		opts = append(opts, websvr.WithReturnOnOutput(true))
	}
	if len(params) > 1 {
		opts = append(opts, websvr.WithResetParams(params[1]))
	}
	return opts
}

// URLFor does another controller handler in this request function.
// it can access any controller method.
func (p *ControllerRegister) URLFor(endpoint string, values ...interface{}) string {
	return (*websvr.ControllerRegister)(p).URLFor(endpoint, values...)
}

// Implement http.Handler interface.
func (p *ControllerRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	(*websvr.ControllerRegister)(p).ServeHTTP(rw, r)
}

// FindRouter Find Router info for URL
func (p *ControllerRegister) FindRouter(ctx *ctxadp.Context) (routerInfo *ControllerInfo, isFind bool) {
	r, ok := (*websvr.ControllerRegister)(p).FindRouter((*context.Context)(ctx))
	return (*ControllerInfo)(r), ok
}

// LogAccess logging info HTTP Access
func LogAccess(ctx *ctxadp.Context, startTime *time.Time, statusCode int) {
	websvr.LogAccess((*context.Context)(ctx), startTime, statusCode)
}
