package new

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
	"os"
	path "path/filepath"
	"strings"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/cmd/utility/commands/version"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/logger/colors"
	"github.com/bhojpur/web/pkg/client/utils"
)

var gopath utils.DocValue
var bhojpurVersion utils.DocValue

var CmdNew = &commands.Command{
	UsageLine: "new [appname] [-gopath=false] [-bhojpur=v1.12.3]",
	Short:     "Creates a Bhojpur Web application",
	Long: `
Creates a Bhojpur Web application for the given app name in the current directory.
  now defaults to generating as a go modules project
  The command 'new' creates a folder named [appname] [-gopath=false] [-bhojpur=v1.12.3] and generates the following structure:

            ├── main.go
            ├── go.mod
            ├── {{"conf"|foldername}}
            │     └── app.conf
            ├── {{"controllers"|foldername}}
            │     └── default.go
            ├── {{"models"|foldername}}
            ├── {{"routers"|foldername}}
            │     └── router.go
            ├── {{"tests"|foldername}}
            │     └── default_test.go
            ├── {{"static"|foldername}}
            │     └── {{"js"|foldername}}
            │     └── {{"css"|foldername}}
            │     └── {{"img"|foldername}}
            └── {{"views"|foldername}}
                  └── index.tpl

`,
	PreRun: nil,
	Run:    CreateApp,
}

var appconf = `appname = {{.Appname}}
httpport = 8080
runmode = dev
`

var maingo = `package main

import (
	_ "{{.Appname}}/routers"
	websvr "github.com/bhojpur/web/pkg/engine"
)

func main() {
	websvr.Run()
}

`
var router = `package routers

import (
	"{{.Appname}}/controllers"
	websvr "github.com/bhojpur/web/pkg/engine"
)

func init() {
    websvr.Router("/", &controllers.MainController{})
}
`
var goMod = `module %s

go %s

require github.com/bhojpur/web %s
require github.com/smartystreets/goconvey v1.6.4
`
var test = `package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"runtime"
	"path/filepath"

    logsvr "github.com/bhojpur/logger/pkg/engine"

	_ "{{.Appname}}/routers"

	websvr "github.com/bhojpur/web/pkg/engine"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".." + string(filepath.Separator))))
	websvr.TestBhojpurInit(apppath)
}


// TestBhojpurWeb is a sample to run an endpoint test
func TestBhojpurWeb(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	websvr.BhojpurApp.Handlers.ServeHTTP(w, r)

	logsvr.Trace("testing", "TestBhojpurWeb", "Code[%d]\n%s", w.Code, w.Body.String())

	Convey("Subject: Test Station Endpoint\n", t, func() {
	        Convey("Status Code Should Be 200", func() {
	                So(w.Code, ShouldEqual, 200)
	        })
	        Convey("The Result Should Not Be Empty", func() {
	                So(w.Body.Len(), ShouldBeGreaterThan, 0)
	        })
	})
}

`

var controllers = `package controllers

import (
	websvr "github.com/bhojpur/web/pkg/engine"
)

type MainController struct {
	websvr.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "bhojpur.net"
	c.Data["Email"] = "info@bhojpur.net"
	c.TplName = "index.tpl"
}
`

var indextpl = `<!DOCTYPE html>

<html>
<head>
  <title>Bhojpur Web - Application</title>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  <link rel="shortcut icon" href="https://static.bhojpur.net/favicon.ico" type="image/x-icon" />

  <style type="text/css">
    *,body {
      margin: 0px;
      padding: 0px;
    }

    body {
      margin: 0px;
      font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
      font-size: 14px;
      line-height: 20px;
      background-color: #fff;
    }

    header,
    footer {
      width: 960px;
      margin-left: auto;
      margin-right: auto;
    }

    .logo {
      background-image: url('https://static.bhojpur.net/image/logo.png');
      background-repeat: no-repeat;
      -webkit-background-size: 100px 100px;
      background-size: 100px 100px;
      background-position: center center;
      text-align: center;
      font-size: 42px;
      padding: 250px 0 70px;
      font-weight: normal;
      text-shadow: 0px 1px 2px #ddd;
    }

    header {
      padding: 100px 0;
    }

    footer {
      line-height: 1.8;
      text-align: center;
      padding: 50px 0;
      color: #999;
    }

    .description {
      text-align: center;
      font-size: 16px;
    }

    a {
      color: #444;
      text-decoration: none;
    }

    .backdrop {
      position: absolute;
      width: 100%;
      height: 100%;
      box-shadow: inset 0px 0px 100px #ddd;
      z-index: -1;
      top: 0px;
      left: 0px;
    }
  </style>
</head>

<body>
  <header>
    <h1 class="logo">Welcome to Bhojpur Web</h1>
    <div class="description">
      The <a href="https://github.com/bhojpur/web">Bhojpur Web</a> is a simple & powerful framework,
	which is inspired by Tornado and Sinatra.
    </div>
  </header>
  <footer>
    <div class="author">
      Official website:
      <a href="https://{{.Website}}">{{.Website}}</a> /
      Contact me:
      <a class="email" href="mailto:{{.Email}}">{{.Email}}</a>
    </div>
  </footer>
  <div class="backdrop"></div>

  <script src="/static/js/reload.min.js"></script>
</body>
</html>
`

var reloadJsClient = `function b(a){var c=new WebSocket(a);c.onclose=function(){setTimeout(function(){b(a)},2E3)};c.onmessage=function(){location.reload()}}try{if(window.WebSocket)try{b("ws://localhost:12450/reload")}catch(a){console.error(a)}else console.log("Your browser does not support WebSockets.")}catch(a){console.error("Exception during connecting to Reload:",a)};
`

func init() {
	CmdNew.Flag.Var(&gopath, "gopath", "Support go path,default false")
	CmdNew.Flag.Var(&bhojpurVersion, "bhojpur", "set Bhojpur Web version,only take effect by go mod")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdNew)
}

func CreateApp(cmd *commands.Command, args []string) int {
	output := cmd.Out()
	if len(args) == 0 {
		cliLogger.Log.Fatal("Argument [appname] is missing")
	}

	if len(args) >= 2 {
		err := cmd.Flag.Parse(args[1:])
		if err != nil {
			cliLogger.Log.Fatal("Parse args err " + err.Error())
		}
	}
	var appPath string
	var packPath string
	var err error
	if gopath == `true` {
		cliLogger.Log.Info("Generate new project support GOPATH")
		version.ShowShortVersionBanner()
		appPath, packPath, err = utils.CheckEnv(args[0])
		if err != nil {
			cliLogger.Log.Fatalf("%s", err)
		}
	} else {
		cliLogger.Log.Info("Generate new project support go modules.")
		appPath = path.Join(utils.GetBhojpurWebWorkPath(), args[0])
		packPath = args[0]
		if bhojpurVersion.String() == `` {
			bhojpurVersion.Set(utils.BHOJPUR_CLI_VERSION)
		}
	}

	if utils.IsExist(appPath) {
		cliLogger.Log.Errorf(colors.Bold("Application '%s' already exists"), appPath)
		cliLogger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes|No] "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	cliLogger.Log.Info("Creating Bhojpur Web application...")

	// If it is the current directory, select the current folder name to package path
	if packPath == "." {
		packPath = path.Base(appPath)
	}

	os.MkdirAll(appPath, 0755)
	if gopath != `true` {
		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "go.mod"), "\x1b[0m")
		utils.WriteToFile(path.Join(appPath, "go.mod"), fmt.Sprintf(goMod, packPath, utils.GetGoVersionSkipMinor(), bhojpurVersion.String()))
	}
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", appPath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "conf"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "conf")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "controllers"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "controllers")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "models"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "models")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "routers"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "routers")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "tests"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "tests")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "static"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "static")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "static", "js"), 0755)
	utils.WriteToFile(path.Join(appPath, "static", "js", "reload.min.js"), reloadJsClient)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "static", "js")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "static", "css"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "static", "css")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "static", "img"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "static", "img")+string(path.Separator), "\x1b[0m")
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "views")+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(appPath, "views"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "conf", "app.conf"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "conf", "app.conf"), strings.Replace(appconf, "{{.Appname}}", path.Base(args[0]), -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "controllers", "default.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "controllers", "default.go"), controllers)

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "views", "index.tpl"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "views", "index.tpl"), indextpl)

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "routers", "router.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "routers", "router.go"), strings.Replace(router, "{{.Appname}}", packPath, -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "tests", "default_test.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "tests", "default_test.go"), strings.Replace(test, "{{.Appname}}", packPath, -1))

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(appPath, "main.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(appPath, "main.go"), strings.Replace(maingo, "{{.Appname}}", packPath, -1))

	cliLogger.Log.Success("New web application successfully created!")
	return 0
}
