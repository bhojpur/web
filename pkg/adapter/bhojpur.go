package adapter

import (
	app "github.com/bhojpur/web/pkg/core"
	web "github.com/bhojpur/web/pkg/engine"
)

const (

	// VERSION represents Bhojpur Web Application framework version.
	VERSION = app.VERSION

	// DEV is for development server engine
	DEV = web.DEV
	// PROD is for production server engine
	PROD = web.PROD
)

// M is Map shortcut
type M web.M

// Hook function to run
type hookfunc func() error

var (
	hooks = make([]hookfunc, 0) // hook function slice to store the hookfunc
)

// AddAPPStartHook is used to register the hookfunc
// The hookfuncs will run in bhojpur.Run()
// such as initiating session , starting middleware , building template, starting admin control and so on.
func AddAPPStartHook(hf ...hookfunc) {
	for _, f := range hf {
		web.AddAPPStartHook(func() error {
			return f()
		})
	}
}

// Run the user's Bhojpur Web application.
// bhojpur.Run() default run on HttpPort
// bhojpur.Run("localhost")
// bhojpur.Run(":8089")
// bhojpur.Run("127.0.0.1:8089")
func Run(params ...string) {
	web.Run(params...)
}

// RunWithMiddleWares Run Bhojpur application with middlewares.
func RunWithMiddleWares(addr string, mws ...MiddleWare) {
	newMws := oldMiddlewareToNew(mws)
	web.RunWithMiddleWares(addr, newMws...)
}

// TestBhojpurInit is for test package init
func TestBhojpurInit(ap string) {
	web.TestBhojpurInit(ap)
}

// InitBhojpurBeforeTest is for test package init
func InitBhojpurBeforeTest(appConfigPath string) {
	web.InitBhojpurBeforeTest(appConfigPath)
}
