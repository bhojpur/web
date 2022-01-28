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
	"os"
	"path/filepath"
	"sync"
)

const (
	// DEV is for developer environment
	DEV = "dev"
	// PROD is for production environment
	PROD = "prod"
)

// M is Map shortcut
type M map[string]interface{}

// Hook function to run
type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) // hook function slice to store the hookfunc
)

// AddAPPStartHook is used to register the hookfunc
// The hookfuncs will run in bhojpur.Run()
// such as initiating session , starting middleware , building template, starting admin control and so on.
func AddAPPStartHook(hf ...hookfunc) {
	hooks = append(hooks, hf...)
}

// Run Bhojpur.NET Platform web server engine.
// WebEngine.Run() default run on HttpPort
// WebEngine.Run("localhost")
// WebEngine.Run(":8089")
// WebEngine.Run("127.0.0.1:8089")
func Run(params ...string) {

	if len(params) > 0 && params[0] != "" {
		WebEngine.Run(params[0])
	}
	WebEngine.Run("")
}

// RunWithMiddleWares Run Bhojpur.NET Platform seb server engine with middlewares.
func RunWithMiddleWares(addr string, mws ...MiddleWare) {
	WebEngine.Run(addr, mws...)
}

var initHttpOnce sync.Once

// TODO move to module init function
func initBeforeHTTPRun() {
	initHttpOnce.Do(func() {
		// init hooks
		AddAPPStartHook(
			registerMime,
			registerDefaultErrorHandler,
			registerSession,
			registerTemplate,
			registerAdmin,
			registerGzip,
			registerCommentRouter,
		)

		for _, hk := range hooks {
			if err := hk(); err != nil {
				panic(err)
			}
		}
	})
}

// TestWebEngineInit is for test package init
func TestWebEngineInit(ap string) {
	path := filepath.Join(ap, "conf", "app.conf")
	os.Chdir(ap)
	InitWebEngineBeforeTest(path)
}

// InitWebEngineBeforeTest is for test package init
func InitWebEngineBeforeTest(appConfigPath string) {
	if err := LoadAppConfig(appConfigProvider, appConfigPath); err != nil {
		panic(err)
	}
	BasConfig.RunMode = "test"
	initBeforeHTTPRun()
}
