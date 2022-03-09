package codegen

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
	"strings"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/pkg/client/codegen/application"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
)

var CmdCodegen = &commands.Command{
	UsageLine: "codegen [command]",
	Short:     "Source code generator",
	Long:      ``,
	Run:       CodeGen,
}

func init() {
	CmdCodegen.Flag.Var(&application.SQL, "sql", "sql file path")
	CmdCodegen.Flag.Var(&application.SQLMode, "sqlmode", "sql mode")
	CmdCodegen.Flag.Var(&application.SQLModePath, "sqlpath", "sql mode path")
	CmdCodegen.Flag.Var(&application.GitRemotePath, "url", "git remote path")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdCodegen)
}

func CodeGen(cmd *commands.Command, args []string) int {
	if len(args) < 1 {
		cliLogger.Log.Fatal("Command is missing")
	}

	if len(args) >= 2 {
		cmd.Flag.Parse(args[1:])
	}

	gcmd := args[0]
	switch gcmd {
	case "gen":
		application.DefaultBhojpurWeb.Run()
	case "toml":
		application.DefaultBhojpurWeb.InitToml()
	case "config":
		application.DefaultBhojpurWeb.GenConfig()
	case "migration":
		application.DefaultBhojpurWeb.Migration(args)
	default:
		cliLogger.Log.Fatal("Command is missing")
	}
	cliLogger.Log.Successf("%s successfully generated!", strings.Title(gcmd))
	return 0
}
