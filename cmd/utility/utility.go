package utility

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
	"github.com/bhojpur/web/cmd/utility/commands"
	_ "github.com/bhojpur/web/cmd/utility/commands/api"
	_ "github.com/bhojpur/web/cmd/utility/commands/bale"
	_ "github.com/bhojpur/web/cmd/utility/commands/codegen"
	_ "github.com/bhojpur/web/cmd/utility/commands/dev"
	_ "github.com/bhojpur/web/cmd/utility/commands/dlv"
	_ "github.com/bhojpur/web/cmd/utility/commands/dockerize"
	_ "github.com/bhojpur/web/cmd/utility/commands/generate"
	_ "github.com/bhojpur/web/cmd/utility/commands/hprose"
	_ "github.com/bhojpur/web/cmd/utility/commands/migrate"
	_ "github.com/bhojpur/web/cmd/utility/commands/new"
	_ "github.com/bhojpur/web/cmd/utility/commands/pack"
	_ "github.com/bhojpur/web/cmd/utility/commands/rs"
	_ "github.com/bhojpur/web/cmd/utility/commands/run"
	_ "github.com/bhojpur/web/cmd/utility/commands/server"
	_ "github.com/bhojpur/web/cmd/utility/commands/update"
	_ "github.com/bhojpur/web/cmd/utility/commands/version"
	"github.com/bhojpur/web/pkg/client/utils"
)

func IfGenerateDocs(name string, args []string) bool {
	if name != "generate" {
		return false
	}
	for _, a := range args {
		if a == "docs" {
			return true
		}
	}
	return false
}

var usageTemplate = `Bhojpur Web CLI utility is a fast and flexible tool for managing your Bhojpur Web application.

You are using CLI utility for Bhojpur Web v2.x. If you are working on Bhojpur Web v1.x, please downgrade version to CLI v1.12.0

{{"USAGE" | headline}}
    {{"webutl command [arguments]" | bold}}

{{"AVAILABLE COMMANDS" | headline}}
{{range .}}{{if .Runnable}}
    {{.Name | printf "%-11s" | bold}} {{.Short}}{{end}}{{end}}

Use {{"webutl help [command]" | bold}} for more information about a command.

{{"ADDITIONAL HELP TOPICS" | headline}}
{{range .}}{{if not .Runnable}}
    {{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use {{"webutl help [topic]" | bold}} for more information about that topic.
`

var helpTemplate = `{{"USAGE" | headline}}
  {{.UsageLine | printf "webutl %s" | bold}}
{{if .Options}}{{endline}}{{"OPTIONS" | headline}}{{range $k,$v := .Options}}
  {{$k | printf "-%s" | bold}}
      {{$v}}
  {{end}}{{end}}
{{"DESCRIPTION" | headline}}
  {{tmpltostr .Long . | trim}}
`

var ErrorTemplate = `bhojpur: %s.
Use {{"webutl help" | bold}} for more information.
`

func Usage() {
	utils.Tmpl(usageTemplate, commands.AvailableCommands)
}

func Help(args []string) {
	if len(args) == 0 {
		Usage()
		return
	}
	if len(args) != 1 {
		utils.PrintErrorAndExit("Too many arguments", ErrorTemplate)
	}

	arg := args[0]

	for _, cmd := range commands.AvailableCommands {
		if cmd.Name() == arg {
			utils.Tmpl(helpTemplate, cmd)
			return
		}
	}
	utils.PrintErrorAndExit("Unknown help topic", ErrorTemplate)
}
