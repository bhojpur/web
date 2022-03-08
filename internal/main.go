//go:generate go run gen/godoc.go
//go:generate go fmt

package main

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

	common "github.com/bhojpur/web/internal/common"
	webapp "github.com/bhojpur/web/pkg/app"
	analytics "github.com/bhojpur/web/pkg/app/analytics"
	webui "github.com/bhojpur/web/pkg/app/ui"
	ctxsvr "github.com/bhojpur/web/pkg/context"
	utils "github.com/bhojpur/web/pkg/core/utils"
	websvr "github.com/bhojpur/web/pkg/engine"
	"github.com/bhojpur/web/pkg/filter/cors"
	"github.com/bhojpur/web/pkg/synthesis"
	test "github.com/bhojpur/web/test"
)

var appengine webapp.Handler

// The main function is the entry point where the Bhojpur web application is
// configured and started. It is executed in two different environments: A
// client (e.g., web browser) and a server (e.g., wasm hosting).
func main() {
	webui.BaseHPadding = 42
	webui.BlockPadding = 18
	analytics.Add(analytics.NewGoogleAnalytics())

	// The first thing to do is to associate the HomePage component with a path.
	// It is done by calling the Route() function, which tells Bhojpur Web what
	// component to display for a given path, on both client and server-side.
	webapp.Route("/", common.NewHomePage())

	// Once the routes are set up, the next thing to do is to either launch the
	// application or the server that serves the application.
	//
	// When executed on the client-side, the webapp.RunWhenOnBrowser() function
	// launches a web application, starting a loop that listens for application
	// events and executes client instructions. Since it is a blocking call, the
	// code below it will never be executed.
	//
	// When executed on the server-side, webapp.RunWhenOnBrowser() does nothing,
	// which lets room for the web server implementation without the need for
	// precompiling instructions.
	webapp.RunWhenOnBrowser()

	websvr.InsertFilter("*", websvr.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// CORS Post method issue
	websvr.InsertFilter("*", websvr.BeforeRouter, func(ctx *ctxsvr.Context) {
		if ctx.Input.Method() == "OPTIONS" {
			ctx.WriteString("ok")
		}
	})

	// websvr.DelStaticPath("/static")
	// websvr.SetStaticPath("/static", "pkg/webui/build/static")

	// Finally, launching the web server that serves the Bhojpur Web application
	// that is done by using the Go standard HTTP package.
	//
	// The Handler is an HTTP handler that serves the client-side application
	// and all its required resources to make it work into a web browser. Here,
	// it is configured to handle requests with a path that starts with "/".
	appengine := webapp.Handler{
		Name:        "Developer's Sandbox",
		Title:       common.DefaultTitle,
		Description: common.DefaultDescription,
		Author:      "Shashi Bhushan Rai",
		Image:       "https://static.bhojpur.net/image/logo.png",
		Keywords: []string{
			"app",
			"webassembly",
			"webapp",
			"web",
			"gui",
			"ui",
			"user interface",
			"frontend",
		},
		BackgroundColor: common.BackgroundColor,
		ThemeColor:      common.BackgroundColor,
		LoadingLabel:    "Bhojpur Web - Developer's Sandbox",
		Scripts: []string{
			"/web/js/prism.js",
		},
		Styles: []string{
			"https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500&display=swap",
			"/web/css/prism.css",
			"/web/css/docs.css",
		},
		RawHeaders:         []string{},
		CacheableResources: []string{},
	} // We must have this for the default web page (i.e. HomePage)
	websvr.AddFuncMap("/", http.HandlerFunc(appengine.ServeHTTP))

	websvr.InsertFilter("/*", websvr.BeforeRouter, StaticContentHandler)
	websvr.Run() // custom configuration read fron ../conf/app.conf file

	websvr.AddFuncMap("*", http.HandlerFunc(appengine.ServeHTTP))
	websvr.AddFuncMap("/अभिवादन/:नाम", http.HandlerFunc(formsHandler))

	// serves static content embedded within the server instance
	websvr.AddFuncMap("/data",
		http.FileServer(
			&synthesis.AssetFS{
				Asset:     test.Asset,
				AssetDir:  test.AssetDir,
				AssetInfo: test.AssetInfo,
				Prefix:    "data",
				Fallback:  "index.html",
			}))
}

// handles web forms
func formsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "नमस्ते, %s!", r.FormValue("नाम"))
}

func StaticContentHandler(ctx *ctxsvr.Context) {
	urlPath := ctx.Request.URL.Path
	path := "."
	if urlPath == "/" {
		path += "/index.html"
	} else {
		path += urlPath
	}

	if utils.FileExists(path) {
		http.ServeFile(ctx.ResponseWriter, ctx.Request, path)
	} else {
		appengine.ServeHTTP(ctx.ResponseWriter, ctx.Request)

	}
}
