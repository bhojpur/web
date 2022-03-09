package hprose

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
	"path"
	"strings"

	"github.com/bhojpur/web/pkg/client/logger/colors"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/cmd/utility/commands/api"
	"github.com/bhojpur/web/cmd/utility/commands/version"
	"github.com/bhojpur/web/pkg/client/generate"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/utils"
)

var CmdHproseapp = &commands.Command{
	// CustomFlags: true,
	UsageLine: "hprose [appname]",
	Short:     "Creates an RPC application based on Hprose and Bhojpur Web frameworks",
	Long: `
  The command 'hprose' creates an RPC application based on both Bhojpur Web and Hprose (http://hprose.com).

  {{"To scaffold out your application, use:"|bold}}

      $ webutl hprose [appname] [-tables=""] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-gopath=false] [-bhojpur=v1.12.3] 

  If 'conn' is empty, the command will generate a sample application. Otherwise the command
  will connect to your database and generate models based on the existing tables.

  The command 'hprose' creates a folder named [appname] with the following structure:

	    ├── main.go
	    ├── go.mod
	    ├── {{"conf"|foldername}}
	    │     └── app.conf
	    └── {{"models"|foldername}}
	          └── object.go
	          └── user.go
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    createhprose,
}

var goMod = `
module %s

go %s

require github.com/bhojpur/web %s
require github.com/smartystreets/goconvey v1.6.4
`

var gopath utils.DocValue
var bhojpurVersion utils.DocValue

func init() {
	CmdHproseapp.Flag.Var(&generate.Tables, "tables", "List of table names separated by a comma.")
	CmdHproseapp.Flag.Var(&generate.SQLDriver, "driver", "Database driver. Either mysql, postgres or sqlite.")
	CmdHproseapp.Flag.Var(&generate.SQLConn, "conn", "Connection string used by the driver to connect to a database instance.")
	CmdHproseapp.Flag.Var(&gopath, "gopath", "Support go path,default false")
	CmdHproseapp.Flag.Var(&bhojpurVersion, "bhojpur", "set Bhojpur Web version, only take effect by go mod")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdHproseapp)
}

func createhprose(cmd *commands.Command, args []string) int {
	output := cmd.Out()
	if len(args) == 0 {
		cliLogger.Log.Fatal("Argument [appname] is missing")
	}

	curpath, _ := os.Getwd()
	if len(args) >= 2 {
		err := cmd.Flag.Parse(args[1:])
		if err != nil {
			cliLogger.Log.Fatal("Parse args err " + err.Error())
		}
	}
	var apppath string
	var packpath string
	var err error
	if gopath == `true` {
		cliLogger.Log.Info("Generate api project support GOPATH")
		version.ShowShortVersionBanner()
		apppath, packpath, err = utils.CheckEnv(args[0])
		if err != nil {
			cliLogger.Log.Fatalf("%s", err)
		}
	} else {
		cliLogger.Log.Info("Generate API project support go modules.")
		apppath = path.Join(utils.GetBhojpurWebWorkPath(), args[0])
		packpath = args[0]
		if bhojpurVersion.String() == `` {
			bhojpurVersion.Set(utils.BHOJPUR_CLI_VERSION)
		}
	}

	if utils.IsExist(apppath) {
		cliLogger.Log.Errorf(colors.Bold("Application '%s' already exists"), apppath)
		cliLogger.Log.Warn(colors.Bold("Do you want to overwrite it? [Yes|No] "))
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	if generate.SQLDriver == "" {
		generate.SQLDriver = "mysql"
	}
	cliLogger.Log.Info("Creating Hprose application...")

	os.MkdirAll(apppath, 0755)
	if gopath != `true` { //generate first for calc model name
		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "go.mod"), "\x1b[0m")
		utils.WriteToFile(path.Join(apppath, "go.mod"), fmt.Sprintf(goMod, packpath, utils.GetGoVersionSkipMinor(), bhojpurVersion.String()))
	}
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath, "\x1b[0m")
	os.Mkdir(path.Join(apppath, "conf"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "conf"), "\x1b[0m")
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "conf", "app.conf"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "conf", "app.conf"),
		strings.Replace(generate.Hproseconf, "{{.Appname}}", args[0], -1))

	if generate.SQLConn != "" {
		cliLogger.Log.Infof("Using '%s' as 'driver'", generate.SQLDriver)
		cliLogger.Log.Infof("Using '%s' as 'conn'", generate.SQLConn)
		cliLogger.Log.Infof("Using '%s' as 'tables'", generate.Tables)
		generate.GenerateHproseAppcode(string(generate.SQLDriver), string(generate.SQLConn), "1", string(generate.Tables), path.Join(curpath, args[0]))

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "main.go"), "\x1b[0m")
		maingoContent := strings.Replace(generate.HproseMainconngo, "{{.Appname}}", packpath, -1)
		maingoContent = strings.Replace(maingoContent, "{{.DriverName}}", string(generate.SQLDriver), -1)
		maingoContent = strings.Replace(maingoContent, "{{HproseFunctionList}}", strings.Join(generate.HproseAddFunctions, ""), -1)
		if generate.SQLDriver == "mysql" {
			maingoContent = strings.Replace(maingoContent, "{{.DriverPkg}}", `_ "github.com/go-sql-driver/mysql"`, -1)
		} else if generate.SQLDriver == "postgres" {
			maingoContent = strings.Replace(maingoContent, "{{.DriverPkg}}", `_ "github.com/lib/pq"`, -1)
		}
		utils.WriteToFile(path.Join(apppath, "main.go"),
			strings.Replace(
				maingoContent,
				"{{.conn}}",
				generate.SQLConn.String(),
				-1,
			),
		)
	} else {
		os.Mkdir(path.Join(apppath, "models"), 0755)
		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "models"), "\x1b[0m")

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "models", "object.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(apppath, "models", "object.go"), api.APIModels)

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "models", "user.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(apppath, "models", "user.go"), api.APIModels2)

		fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "main.go"), "\x1b[0m")
		utils.WriteToFile(path.Join(apppath, "main.go"),
			strings.Replace(generate.HproseMaingo, "{{.Appname}}", packpath, -1))
	}
	cliLogger.Log.Success("New Hprose application successfully created!")
	return 0
}
