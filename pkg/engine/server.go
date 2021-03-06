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
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"
	"time"

	"golang.org/x/crypto/acme/autocert"

	logsvr "github.com/bhojpur/logger/pkg/engine"
	ctxsvr "github.com/bhojpur/web/pkg/context"
	"github.com/bhojpur/web/pkg/core/utils"
	"github.com/bhojpur/web/pkg/grace"
)

// BhojpurApp is an application instance
// If you are using single server, you could use this
// But if you need multiple servers, do not use this
var BhojpurApp *HttpServer

func init() {
	// create Bhojpur.NET Platform application
	BhojpurApp = NewHttpSever()
}

// HttpServer defines Bhojpur.NET Platform application with a new PatternServeMux.
type HttpServer struct {
	Handlers *ControllerRegister
	Server   *http.Server
	Cfg      *Config
}

// NewHttpSever returns a new Bhojpur.NET Platform application.
// this method will use the BConfig as the configure to create HttpServer
// Be careful that when you update BConfig, the server's Cfg will be updated too
func NewHttpSever() *HttpServer {
	return NewHttpServerWithCfg(BConfig)
}

// NewHttpServerWithCfg will create an sever with specific cfg
func NewHttpServerWithCfg(cfg *Config) *HttpServer {
	cr := NewControllerRegisterWithCfg(cfg)
	app := &HttpServer{
		Handlers: cr,
		Server:   &http.Server{},
		Cfg:      cfg,
	}

	return app
}

// MiddleWare function for http.Handler
type MiddleWare func(http.Handler) http.Handler

// Run the Bhojpur.NET Platform's primary web server engine.
func (app *HttpServer) Run(addr string, mws ...MiddleWare) {
	initBeforeHTTPRun()

	// init...
	app.initAddr(addr)
	app.Handlers.Init()

	addr = app.Cfg.Listen.HTTPAddr

	if app.Cfg.Listen.HTTPPort != 0 {
		addr = fmt.Sprintf("%s:%d", app.Cfg.Listen.HTTPAddr, app.Cfg.Listen.HTTPPort)
	}

	var (
		err        error
		l          net.Listener
		endRunning = make(chan bool, 1)
	)

	// run cgi server
	if app.Cfg.Listen.EnableFcgi {
		if app.Cfg.Listen.EnableStdIo {
			if err = fcgi.Serve(nil, app.Handlers); err == nil { // standard I/O
				logsvr.Info("Use Fast CGI via standard I/O")
			} else {
				logsvr.Critical("Cannot use Fast CGI via standard I/O", err)
			}
			return
		}
		if app.Cfg.Listen.HTTPPort == 0 {
			// remove the Socket file before start
			if utils.FileExists(addr) {
				os.Remove(addr)
			}
			l, err = net.Listen("unix", addr)
		} else {
			l, err = net.Listen("tcp", addr)
		}
		if err != nil {
			logsvr.Critical("Listen: ", err)
		}
		if err = fcgi.Serve(l, app.Handlers); err != nil {
			logsvr.Critical("fcgi.Serve: ", err)
		}
		return
	}

	app.Server.Handler = app.Handlers
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] == nil {
			continue
		}
		app.Server.Handler = mws[i](app.Server.Handler)
	}
	app.Server.ReadTimeout = time.Duration(app.Cfg.Listen.ServerTimeOut) * time.Second
	app.Server.WriteTimeout = time.Duration(app.Cfg.Listen.ServerTimeOut) * time.Second
	app.Server.ErrorLog = logsvr.GetLogger("HTTP")

	// run graceful mode
	if app.Cfg.Listen.Graceful {
		httpsAddr := app.Cfg.Listen.HTTPSAddr
		app.Server.Addr = httpsAddr
		if app.Cfg.Listen.EnableHTTPS || app.Cfg.Listen.EnableMutualHTTPS {
			go func() {
				time.Sleep(1000 * time.Microsecond)
				if app.Cfg.Listen.HTTPSPort != 0 {
					httpsAddr = fmt.Sprintf("%s:%d", app.Cfg.Listen.HTTPSAddr, app.Cfg.Listen.HTTPSPort)
					app.Server.Addr = httpsAddr
				}
				server := grace.NewServer(httpsAddr, app.Server.Handler)
				server.Server.ReadTimeout = app.Server.ReadTimeout
				server.Server.WriteTimeout = app.Server.WriteTimeout
				if app.Cfg.Listen.EnableMutualHTTPS {
					if err := server.ListenAndServeMutualTLS(app.Cfg.Listen.HTTPSCertFile,
						app.Cfg.Listen.HTTPSKeyFile,
						app.Cfg.Listen.TrustCaFile); err != nil {
						logsvr.Critical("ListenAndServeTLS: ", err, fmt.Sprintf("%d", os.Getpid()))
						time.Sleep(100 * time.Microsecond)
					}
				} else {
					if app.Cfg.Listen.AutoTLS {
						m := autocert.Manager{
							Prompt:     autocert.AcceptTOS,
							HostPolicy: autocert.HostWhitelist(app.Cfg.Listen.Domains...),
							Cache:      autocert.DirCache(app.Cfg.Listen.TLSCacheDir),
						}
						app.Server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
						app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile = "", ""
					}
					if err := server.ListenAndServeTLS(app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile); err != nil {
						logsvr.Critical("ListenAndServeTLS: ", err, fmt.Sprintf("%d", os.Getpid()))
						time.Sleep(100 * time.Microsecond)
					}
				}
				endRunning <- true
			}()
		}
		if app.Cfg.Listen.EnableHTTP {
			go func() {
				server := grace.NewServer(addr, app.Server.Handler)
				server.Server.ReadTimeout = app.Server.ReadTimeout
				server.Server.WriteTimeout = app.Server.WriteTimeout
				if app.Cfg.Listen.ListenTCP4 {
					server.Network = "tcp4"
				}
				if err := server.ListenAndServe(); err != nil {
					logsvr.Critical("ListenAndServe: ", err, fmt.Sprintf("%d", os.Getpid()))
					time.Sleep(100 * time.Microsecond)
				}
				endRunning <- true
			}()
		}
		<-endRunning
		return
	}

	// run normal mode
	if app.Cfg.Listen.EnableHTTPS || app.Cfg.Listen.EnableMutualHTTPS {
		go func() {
			time.Sleep(1000 * time.Microsecond)
			if app.Cfg.Listen.HTTPSPort != 0 {
				app.Server.Addr = fmt.Sprintf("%s:%d", app.Cfg.Listen.HTTPSAddr, app.Cfg.Listen.HTTPSPort)
			} else if app.Cfg.Listen.EnableHTTP {
				logsvr.Info("Start https server error, conflict with http. Please reset https port")
				return
			}
			logsvr.Info("Bhojpur WebEngine - HTTPS server running on https://%s", app.Server.Addr)
			if app.Cfg.Listen.AutoTLS {
				m := autocert.Manager{
					Prompt:     autocert.AcceptTOS,
					HostPolicy: autocert.HostWhitelist(app.Cfg.Listen.Domains...),
					Cache:      autocert.DirCache(app.Cfg.Listen.TLSCacheDir),
				}
				app.Server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
				app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile = "", ""
			} else if app.Cfg.Listen.EnableMutualHTTPS {
				pool := x509.NewCertPool()
				data, err := ioutil.ReadFile(app.Cfg.Listen.TrustCaFile)
				if err != nil {
					logsvr.Info("MutualHTTPS should provide TrustCaFile")
					return
				}
				pool.AppendCertsFromPEM(data)
				app.Server.TLSConfig = &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.ClientAuthType(app.Cfg.Listen.ClientAuth),
				}
			}
			if err := app.Server.ListenAndServeTLS(app.Cfg.Listen.HTTPSCertFile, app.Cfg.Listen.HTTPSKeyFile); err != nil {
				logsvr.Critical("ListenAndServeTLS: ", err)
				time.Sleep(100 * time.Microsecond)
				endRunning <- true
			}
		}()
	}
	if app.Cfg.Listen.EnableHTTP {
		go func() {
			app.Server.Addr = addr
			logsvr.Info("Bhojpur WebEngine - HTTP server running on http://%s", app.Server.Addr)
			if app.Cfg.Listen.ListenTCP4 {
				ln, err := net.Listen("tcp4", app.Server.Addr)
				if err != nil {
					logsvr.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
					return
				}
				if err = app.Server.Serve(ln); err != nil {
					logsvr.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
					return
				}
			} else {
				if err := app.Server.ListenAndServe(); err != nil {
					logsvr.Critical("ListenAndServe: ", err)
					time.Sleep(100 * time.Microsecond)
					endRunning <- true
				}
			}
		}()
	}
	<-endRunning
}

// Router see HttpServer.Router
func Router(rootpath string, c ControllerInterface, mappingMethods ...string) *HttpServer {
	return RouterWithOpts(rootpath, c, WithRouterMethods(c, mappingMethods...))
}

func RouterWithOpts(rootpath string, c ControllerInterface, opts ...ControllerOption) *HttpServer {
	return BhojpurApp.RouterWithOpts(rootpath, c, opts...)
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
func (app *HttpServer) Router(rootPath string, c ControllerInterface, mappingMethods ...string) *HttpServer {
	return app.RouterWithOpts(rootPath, c, WithRouterMethods(c, mappingMethods...))
}

func (app *HttpServer) RouterWithOpts(rootPath string, c ControllerInterface, opts ...ControllerOption) *HttpServer {
	app.Handlers.Add(rootPath, c, opts...)
	return app
}

// UnregisterFixedRoute see HttpServer.UnregisterFixedRoute
func UnregisterFixedRoute(fixedRoute string, method string) *HttpServer {
	return BhojpurApp.UnregisterFixedRoute(fixedRoute, method)
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
func (app *HttpServer) UnregisterFixedRoute(fixedRoute string, method string) *HttpServer {
	subPaths := splitPath(fixedRoute)
	if method == "" || method == "*" {
		for m := range HTTPMETHOD {
			if _, ok := app.Handlers.routers[m]; !ok {
				continue
			}
			if app.Handlers.routers[m].prefix == strings.Trim(fixedRoute, "/ ") {
				findAndRemoveSingleTree(app.Handlers.routers[m])
				continue
			}
			findAndRemoveTree(subPaths, app.Handlers.routers[m], m)
		}
		return app
	}
	// Single HTTP method
	um := strings.ToUpper(method)
	if _, ok := app.Handlers.routers[um]; ok {
		if app.Handlers.routers[um].prefix == strings.Trim(fixedRoute, "/ ") {
			findAndRemoveSingleTree(app.Handlers.routers[um])
			return app
		}
		findAndRemoveTree(subPaths, app.Handlers.routers[um], um)
	}
	return app
}

func findAndRemoveTree(paths []string, entryPointTree *Tree, method string) {
	for i := range entryPointTree.fixrouters {
		if entryPointTree.fixrouters[i].prefix == paths[0] {
			if len(paths) == 1 {
				if len(entryPointTree.fixrouters[i].fixrouters) > 0 {
					// If the route had children subtrees, remove just the functional leaf,
					// to allow children to function as before
					if len(entryPointTree.fixrouters[i].leaves) > 0 {
						entryPointTree.fixrouters[i].leaves[0] = nil
						entryPointTree.fixrouters[i].leaves = entryPointTree.fixrouters[i].leaves[1:]
					}
				} else {
					// Remove the *Tree from the fixrouters slice
					entryPointTree.fixrouters[i] = nil

					if i == len(entryPointTree.fixrouters)-1 {
						entryPointTree.fixrouters = entryPointTree.fixrouters[:i]
					} else {
						entryPointTree.fixrouters = append(entryPointTree.fixrouters[:i], entryPointTree.fixrouters[i+1:len(entryPointTree.fixrouters)]...)
					}
				}
				return
			}
			findAndRemoveTree(paths[1:], entryPointTree.fixrouters[i], method)
		}
	}
}

func findAndRemoveSingleTree(entryPointTree *Tree) {
	if entryPointTree == nil {
		return
	}
	if len(entryPointTree.fixrouters) > 0 {
		// If the route had children subtrees, remove just the functional leaf,
		// to allow children to function as before
		if len(entryPointTree.leaves) > 0 {
			entryPointTree.leaves[0] = nil
			entryPointTree.leaves = entryPointTree.leaves[1:]
		}
	}
}

// Include see HttpServer.Include
func Include(cList ...ControllerInterface) *HttpServer {
	return BhojpurApp.Include(cList...)
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
func (app *HttpServer) Include(cList ...ControllerInterface) *HttpServer {
	app.Handlers.Include(cList...)
	return app
}

// RESTRouter see HttpServer.RESTRouter
func RESTRouter(rootpath string, c ControllerInterface) *HttpServer {
	return BhojpurApp.RESTRouter(rootpath, c)
}

// RESTRouter adds a restful controller handler to BhojpurApp.
// its' controller implements websvr.ControllerInterface and
// defines a param "pattern/:objectId" to visit each resource.
func (app *HttpServer) RESTRouter(rootpath string, c ControllerInterface) *HttpServer {
	app.Router(rootpath, c)
	app.Router(path.Join(rootpath, ":objectId"), c)
	return app
}

// AutoRouter see HttpServer.AutoRouter
func AutoRouter(c ControllerInterface) *HttpServer {
	return BhojpurApp.AutoRouter(c)
}

// AutoRouter adds defined controller handler to BhojpurApp.
// it's same to HttpServer.AutoRouter.
// if websvr.AddAuto(&MainController{}) and MainController has methods List and Page,
// visit the url /main/list to exec List function or /main/page to exec Page function.
func (app *HttpServer) AutoRouter(c ControllerInterface) *HttpServer {
	app.Handlers.AddAuto(c)
	return app
}

// AutoPrefix see HttpServer.AutoPrefix
func AutoPrefix(prefix string, c ControllerInterface) *HttpServer {
	return BhojpurApp.AutoPrefix(prefix, c)
}

// AutoPrefix adds controller handler to BhojpurApp with prefix.
// it's same to HttpServer.AutoRouterWithPrefix.
// if websvr.AutoPrefix("/admin",&MainController{}) and MainController has methods List and Page,
// visit the url /admin/main/list to exec List function or /admin/main/page to exec Page function.
func (app *HttpServer) AutoPrefix(prefix string, c ControllerInterface) *HttpServer {
	app.Handlers.AddAutoPrefix(prefix, c)
	return app
}

// CtrlGet see HttpServer.CtrlGet
func CtrlGet(rootpath string, f interface{}) {
	BhojpurApp.CtrlGet(rootpath, f)
}

// CtrlGet used to register router for CtrlGet method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlGet("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlGet(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlGet(rootpath, f)
	return app
}

// CtrlPost see HttpServer.CtrlGet
func CtrlPost(rootpath string, f interface{}) {
	BhojpurApp.CtrlPost(rootpath, f)
}

// CtrlPost used to register router for CtrlPost method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlPost("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlPost(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlPost(rootpath, f)
	return app
}

// CtrlHead see HttpServer.CtrlHead
func CtrlHead(rootpath string, f interface{}) {
	BhojpurApp.CtrlHead(rootpath, f)
}

// CtrlHead used to register router for CtrlHead method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlHead("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlHead(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlHead(rootpath, f)
	return app
}

// CtrlPut see HttpServer.CtrlPut
func CtrlPut(rootpath string, f interface{}) {
	BhojpurApp.CtrlPut(rootpath, f)
}

// CtrlPut used to register router for CtrlPut method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlPut("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlPut(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlPut(rootpath, f)
	return app
}

// CtrlPatch see HttpServer.CtrlPatch
func CtrlPatch(rootpath string, f interface{}) {
	BhojpurApp.CtrlPatch(rootpath, f)
}

// CtrlPatch used to register router for CtrlPatch method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlPatch("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlPatch(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlPatch(rootpath, f)
	return app
}

// CtrlDelete see HttpServer.CtrlDelete
func CtrlDelete(rootpath string, f interface{}) {
	BhojpurApp.CtrlDelete(rootpath, f)
}

// CtrlDelete used to register router for CtrlDelete method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlDelete("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlDelete(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlDelete(rootpath, f)
	return app
}

// CtrlOptions see HttpServer.CtrlOptions
func CtrlOptions(rootpath string, f interface{}) {
	BhojpurApp.CtrlOptions(rootpath, f)
}

// CtrlOptions used to register router for CtrlOptions method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlOptions("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlOptions(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlOptions(rootpath, f)
	return app
}

// CtrlAny see HttpServer.CtrlAny
func CtrlAny(rootpath string, f interface{}) {
	BhojpurApp.CtrlAny(rootpath, f)
}

// CtrlAny used to register router for CtrlAny method
// usage:
//    type MyController struct {
//	     web.Controller
//    }
//    func (m MyController) Ping() {
//	     m.Ctx.Output.Body([]byte("hello world"))
//    }
//
//    CtrlAny("/api/:id", MyController.Ping)
func (app *HttpServer) CtrlAny(rootpath string, f interface{}) *HttpServer {
	app.Handlers.CtrlAny(rootpath, f)
	return app
}

// Get see HttpServer.Get
func Get(rootpath string, f HandleFunc) *HttpServer {
	return BhojpurApp.Get(rootpath, f)
}

// Get used to register router for Get method
// usage:
//    websvr.Get("/", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Get(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Get(rootpath, f)
	return app
}

// Post see HttpServer.Post
func Post(rootpath string, f HandleFunc) *HttpServer {
	return BhojpurApp.Post(rootpath, f)
}

// Post used to register router for Post method
// usage:
//    websvr.Post("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Post(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Post(rootpath, f)
	return app
}

// Delete see HttpServer.Delete
func Delete(rootpath string, f HandleFunc) *HttpServer {
	return BhojpurApp.Delete(rootpath, f)
}

// Delete used to register router for Delete method
// usage:
//    websvr.Delete("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Delete(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Delete(rootpath, f)
	return app
}

// Put see HttpServer.Put
func Put(rootpath string, f HandleFunc) *HttpServer {
	return BhojpurApp.Put(rootpath, f)
}

// Put used to register router for Put method
// usage:
//    websvr.Put("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Put(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Put(rootpath, f)
	return app
}

// Head see HttpServer.Head
func Head(rootpath string, f HandleFunc) *HttpServer {
	return BhojpurApp.Head(rootpath, f)
}

// Head used to register router for Head method
// usage:
//    websvr.Head("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Head(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Head(rootpath, f)
	return app
}

// Options see HttpServer.Options
func Options(rootpath string, f HandleFunc) *HttpServer {
	BhojpurApp.Handlers.Options(rootpath, f)
	return BhojpurApp
}

// Options used to register router for Options method
// usage:
//    websvr.Options("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Options(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Options(rootpath, f)
	return app
}

// Patch see HttpServer.Patch
func Patch(rootpath string, f HandleFunc) *HttpServer {
	return BhojpurApp.Patch(rootpath, f)
}

// Patch used to register router for Patch method
// usage:
//    websvr.Patch("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Patch(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Patch(rootpath, f)
	return app
}

// Any see HttpServer.Any
func Any(rootpath string, f HandleFunc) *HttpServer {
	return BhojpurApp.Any(rootpath, f)
}

// Any used to register router for all methods
// usage:
//    websvr.Any("/api", func(ctx *context.Context){
//          ctx.Output.Body("hello world")
//    })
func (app *HttpServer) Any(rootpath string, f HandleFunc) *HttpServer {
	app.Handlers.Any(rootpath, f)
	return app
}

// Handler see HttpServer.Handler
func Handler(rootpath string, h http.Handler, options ...interface{}) *HttpServer {
	return BhojpurApp.Handler(rootpath, h, options...)
}

// Handler used to register a Handler router
// usage:
//    websvr.Handler("/api", http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
//          fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
//    }))
func (app *HttpServer) Handler(rootpath string, h http.Handler, options ...interface{}) *HttpServer {
	app.Handlers.Handler(rootpath, h, options...)
	return app
}

// InserFilter see HttpServer.InsertFilter
func InsertFilter(pattern string, pos int, filter FilterFunc, opts ...FilterOpt) *HttpServer {
	return BhojpurApp.InsertFilter(pattern, pos, filter, opts...)
}

// InsertFilter adds a FilterFunc with pattern condition and action constant.
// The pos means action constant including
// websvr.BeforeStatic, websvr.BeforeRouter, websvr.BeforeExec, websvr.AfterExec and websvr.FinishRouter.
// The bool params is for setting the returnOnOutput value (false allows multiple filters to execute)
func (app *HttpServer) InsertFilter(pattern string, pos int, filter FilterFunc, opts ...FilterOpt) *HttpServer {
	app.Handlers.InsertFilter(pattern, pos, filter, opts...)
	return app
}

// InsertFilterChain see HttpServer.InsertFilterChain
func InsertFilterChain(pattern string, filterChain FilterChain, opts ...FilterOpt) *HttpServer {
	return BhojpurApp.InsertFilterChain(pattern, filterChain, opts...)
}

// InsertFilterChain adds a FilterFunc built by filterChain.
// This filter will be executed before all filters.
// the filter's behavior like stack's behavior
// and the last filter is serving the http request
func (app *HttpServer) InsertFilterChain(pattern string, filterChain FilterChain, opts ...FilterOpt) *HttpServer {
	app.Handlers.InsertFilterChain(pattern, filterChain, opts...)
	return app
}

func (app *HttpServer) initAddr(addr string) {
	strs := strings.Split(addr, ":")
	if len(strs) > 0 && strs[0] != "" {
		app.Cfg.Listen.HTTPAddr = strs[0]
		app.Cfg.Listen.Domains = []string{strs[0]}
	}
	if len(strs) > 1 && strs[1] != "" {
		app.Cfg.Listen.HTTPPort, _ = strconv.Atoi(strs[1])
	}
}

func (app *HttpServer) LogAccess(ctx *ctxsvr.Context, startTime *time.Time, statusCode int) {
	// Skip logging if AccessLogs config is false
	if !app.Cfg.Log.AccessLogs {
		return
	}
	// Skip logging static requests unless EnableStaticLogs config is true
	if !app.Cfg.Log.EnableStaticLogs && DefaultAccessLogFilter.Filter(ctx) {
		return
	}
	var (
		requestTime time.Time
		elapsedTime time.Duration
		r           = ctx.Request
	)
	if startTime != nil {
		requestTime = *startTime
		elapsedTime = time.Since(*startTime)
	}
	record := &logsvr.AccessLogRecord{
		RemoteAddr:     ctx.Input.IP(),
		RequestTime:    requestTime,
		RequestMethod:  r.Method,
		Request:        fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		ServerProtocol: r.Proto,
		Host:           r.Host,
		Status:         statusCode,
		ElapsedTime:    elapsedTime,
		HTTPReferrer:   r.Header.Get("Referer"),
		HTTPUserAgent:  r.Header.Get("User-Agent"),
		RemoteUser:     r.Header.Get("Remote-User"),
		BodyBytesSent:  r.ContentLength,
	}
	logsvr.AccessLog(record, app.Cfg.Log.AccessLogsFormat)
}

// PrintTree prints all registered routers.
func (app *HttpServer) PrintTree() M {
	var (
		content     = M{}
		methods     = []string{}
		methodsData = make(M)
	)
	for method, t := range app.Handlers.routers {

		resultList := new([][]string)

		printTree(resultList, t)

		methods = append(methods, template.HTMLEscapeString(method))
		methodsData[template.HTMLEscapeString(method)] = resultList
	}

	content["Data"] = methodsData
	content["Methods"] = methods
	return content
}

func printTree(resultList *[][]string, t *Tree) {
	for _, tr := range t.fixrouters {
		printTree(resultList, tr)
	}
	if t.wildcard != nil {
		printTree(resultList, t.wildcard)
	}
	for _, l := range t.leaves {
		if v, ok := l.runObject.(*ControllerInfo); ok {
			if v.routerType == routerTypeBhojpur {
				result := []string{
					template.HTMLEscapeString(v.pattern),
					template.HTMLEscapeString(fmt.Sprintf("%s", v.methods)),
					template.HTMLEscapeString(v.controllerType.String()),
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeRESTFul {
				result := []string{
					template.HTMLEscapeString(v.pattern),
					template.HTMLEscapeString(fmt.Sprintf("%s", v.methods)),
					"",
				}
				*resultList = append(*resultList, result)
			} else if v.routerType == routerTypeHandler {
				result := []string{
					template.HTMLEscapeString(v.pattern),
					"",
					"",
				}
				*resultList = append(*resultList, result)
			}
		}
	}
}

func (app *HttpServer) reportFilter() M {
	filterTypeData := make(M)
	// filterTypes := []string{}
	if app.Handlers.enableFilter {
		// var filterType string
		for k, fr := range map[int]string{
			BeforeStatic: "Before Static",
			BeforeRouter: "Before Router",
			BeforeExec:   "Before Exec",
			AfterExec:    "After Exec",
			FinishRouter: "Finish Router",
		} {
			if bf := app.Handlers.filters[k]; len(bf) > 0 {
				resultList := new([][]string)
				for _, f := range bf {
					result := []string{
						// void xss
						template.HTMLEscapeString(f.pattern),
						template.HTMLEscapeString(utils.GetFuncName(f.filterFunc)),
					}
					*resultList = append(*resultList, result)
				}
				filterTypeData[fr] = resultList
			}
		}
	}

	return filterTypeData
}
