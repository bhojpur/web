package generate

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
	"strings"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/cmd/utility/commands/version"
	"github.com/bhojpur/web/pkg/client/config"
	"github.com/bhojpur/web/pkg/client/generate"
	"github.com/bhojpur/web/pkg/client/generate/swaggergen"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/utils"
)

var CmdGenerate = &commands.Command{
	UsageLine: "generate [command]",
	Short:     "Source code generator",
	Long: `▶ {{"To scaffold out your entire application:"|bold}}

     $ webutl generate scaffold [scaffoldname] [-fields="title:string,body:text"] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"]

  ▶ {{"To generate a Model based on fields:"|bold}}

     $ webutl generate model [modelname] [-fields="name:type"]

  ▶ {{"To generate a controller:"|bold}}

     $ webutl generate controller [controllerfile]

  ▶ {{"To generate a CRUD view:"|bold}}

     $ webutl generate view [viewpath]

  ▶ {{"To generate a migration file for making database schema updates:"|bold}}

     $ webutl generate migration [migrationfile] [-fields="name:type"]

  ▶ {{"To generate swagger doc file:"|bold}}

     $ webutl generate docs

    ▶ {{"To generate swagger doc file:"|bold}}

     $ webutl generate routers [-ctrlDir=/path/to/controller/directory] [-routersFile=/path/to/routers/file.go] [-routersPkg=myPackage]

  ▶ {{"To generate a test case:"|bold}}

     $ webutl generate test [routerfile]

  ▶ {{"To generate appcode based on an existing database:"|bold}}

     $ webutl generate appcode [-tables=""] [-driver=mysql] [-conn="root:@tcp(127.0.0.1:3306)/test"] [-level=3]
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    GenerateCode,
}

func init() {
	CmdGenerate.Flag.Var(&generate.Tables, "tables", "List of table names separated by a comma.")
	CmdGenerate.Flag.Var(&generate.SQLDriver, "driver", "Database SQLDriver. Either mysql, postgres or sqlite.")
	CmdGenerate.Flag.Var(&generate.SQLConn, "conn", "Connection string used by the SQLDriver to connect to a database instance.")
	CmdGenerate.Flag.Var(&generate.Level, "level", "Either 1, 2 or 3. i.e. 1=models; 2=models and controllers; 3=models, controllers and routers.")
	CmdGenerate.Flag.Var(&generate.Fields, "fields", "List of table Fields.")
	CmdGenerate.Flag.Var(&generate.DDL, "ddl", "Generate DDL Migration")

	// Bhojpur Web generate routers
	CmdGenerate.Flag.Var(&generate.ControllerDirectory, "ctrlDir",
		"Controller directory. Bhojpur Web scans this directory and its sub directory to generate routers")
	CmdGenerate.Flag.Var(&generate.RoutersFile, "routersFile",
		"Routers file. If not found, Bhojpur Web create a new one. Bhojpur Web will truncates this file and output routers info into this file")
	CmdGenerate.Flag.Var(&generate.RouterPkg, "routersPkg",
		`router's package. Default is routers, it means that "package routers" in the generated file`)

	commands.AvailableCommands = append(commands.AvailableCommands, CmdGenerate)
}

func GenerateCode(cmd *commands.Command, args []string) int {
	currpath, _ := os.Getwd()
	if len(args) < 1 {
		cliLogger.Log.Fatal("Command is missing")
	}

	gcmd := args[0]
	switch gcmd {
	case "scaffold":
		scaffold(cmd, args, currpath)
	case "docs":
		swaggergen.GenerateDocs(currpath)
	case "appcode":
		appCode(cmd, args, currpath)
	case "migration":
		migration(cmd, args, currpath)
	case "controller":
		controller(args, currpath)
	case "model":
		model(cmd, args, currpath)
	case "view":
		view(args, currpath)
	case "routers":
		genRouters(cmd, args)
	default:
		cliLogger.Log.Fatal("Command is missing")
	}
	cliLogger.Log.Successf("%s successfully generated!", strings.Title(gcmd))
	return 0
}

func genRouters(cmd *commands.Command, args []string) {
	err := cmd.Flag.Parse(args[1:])
	cliLogger.Log.Infof("input parameter: %v", args)
	if err != nil {
		cliLogger.Log.Errorf("could not parse input parameter: %+v", err)
		return
	}
	generate.GenRouters()
}

func scaffold(cmd *commands.Command, args []string, currpath string) {
	if len(args) < 2 {
		cliLogger.Log.Fatal("Wrong number of arguments. Run: webutl help generate")
	}

	cmd.Flag.Parse(args[2:])
	if generate.SQLDriver == "" {
		generate.SQLDriver = utils.DocValue(config.Conf.Database.Driver)
		if generate.SQLDriver == "" {
			generate.SQLDriver = "mysql"
		}
	}
	if generate.SQLConn == "" {
		generate.SQLConn = utils.DocValue(config.Conf.Database.Conn)
		if generate.SQLConn == "" {
			generate.SQLConn = "root:@tcp(127.0.0.1:3306)/test"
		}
	}
	if generate.Fields == "" {
		cliLogger.Log.Hint("Fields option should not be empty, i.e. -Fields=\"title:string,body:text\"")
		cliLogger.Log.Fatal("Wrong number of arguments. Run: webutl help generate")
	}
	sname := args[1]
	generate.GenerateScaffold(sname, generate.Fields.String(), currpath, generate.SQLDriver.String(), generate.SQLConn.String())
}

func appCode(cmd *commands.Command, args []string, currpath string) {
	cmd.Flag.Parse(args[1:])
	if generate.SQLDriver == "" {
		generate.SQLDriver = utils.DocValue(config.Conf.Database.Driver)
		if generate.SQLDriver == "" {
			generate.SQLDriver = "mysql"
		}
	}
	if generate.SQLConn == "" {
		generate.SQLConn = utils.DocValue(config.Conf.Database.Conn)
		if generate.SQLConn == "" {
			if generate.SQLDriver == "mysql" {
				generate.SQLConn = "root:@tcp(127.0.0.1:3306)/test"
			} else if generate.SQLDriver == "postgres" {
				generate.SQLConn = "postgres://postgres:postgres@127.0.0.1:5432/postgres"
			}
		}
	}
	if generate.Level == "" {
		generate.Level = "3"
	}
	cliLogger.Log.Infof("Using '%s' as 'SQLDriver'", generate.SQLDriver)
	cliLogger.Log.Infof("Using '%s' as 'SQLConn'", generate.SQLConn)
	cliLogger.Log.Infof("Using '%s' as 'Tables'", generate.Tables)
	cliLogger.Log.Infof("Using '%s' as 'Level'", generate.Level)
	generate.GenerateAppcode(generate.SQLDriver.String(), generate.SQLConn.String(), generate.Level.String(), generate.Tables.String(), currpath)
}

func migration(cmd *commands.Command, args []string, currpath string) {
	if len(args) < 2 {
		cliLogger.Log.Fatal("Wrong number of arguments. Run: webutl help generate")
	}
	cmd.Flag.Parse(args[2:])
	mname := args[1]

	cliLogger.Log.Infof("Using '%s' as migration name", mname)

	upsql := ""
	downsql := ""
	if generate.Fields != "" {
		dbMigrator := generate.NewDBDriver()
		upsql = dbMigrator.GenerateCreateUp(mname)
		downsql = dbMigrator.GenerateCreateDown(mname)
	}
	generate.GenerateMigration(mname, upsql, downsql, currpath)
}

func controller(args []string, currpath string) {
	if len(args) == 2 {
		cname := args[1]
		generate.GenerateController(cname, currpath)
	} else {
		cliLogger.Log.Fatal("Wrong number of arguments. Run: webutl help generate")
	}
}

func model(cmd *commands.Command, args []string, currpath string) {
	if len(args) < 2 {
		cliLogger.Log.Fatal("Wrong number of arguments. Run: webutl help generate")
	}
	cmd.Flag.Parse(args[2:])
	if generate.Fields == "" {
		cliLogger.Log.Hint("Fields option should not be empty, i.e. -Fields=\"title:string,body:text\"")
		cliLogger.Log.Fatal("Wrong number of arguments. Run: webutl help generate")
	}
	sname := args[1]
	generate.GenerateModel(sname, generate.Fields.String(), currpath)
}

func view(args []string, currpath string) {
	if len(args) == 2 {
		cname := args[1]
		generate.GenerateView(cname, currpath)
	} else {
		cliLogger.Log.Fatal("Wrong number of arguments. Run: webutl help generate")
	}
}
