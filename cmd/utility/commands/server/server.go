package apiapp

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

	cliLogger "github.com/bhojpur/web/pkg/client/logger"

	"os"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/cmd/utility/commands/version"
	"github.com/bhojpur/web/pkg/client/utils"
)

var CmdServer = &commands.Command{
	// CustomFlags: true,
	UsageLine: "server [port]",
	Short:     "serving static content over HTTP on port",
	Long: `
  The command 'server' creates a Bhojpur Web API application.
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    createAPI,
}

var (
	a utils.DocValue
	p utils.DocValue
	f utils.DocValue
)

func init() {
	CmdServer.Flag.Var(&a, "a", "Listen address")
	CmdServer.Flag.Var(&p, "p", "Listen port")
	CmdServer.Flag.Var(&f, "f", "Static files fold")
	commands.AvailableCommands = append(commands.AvailableCommands, CmdServer)
}

func createAPI(cmd *commands.Command, args []string) int {
	if len(args) > 0 {
		err := cmd.Flag.Parse(args[1:])
		if err != nil {
			cliLogger.Log.Error(err.Error())
		}
	}
	if a == "" {
		a = "127.0.0.1"
	}
	if p == "" {
		p = "8080"
	}
	if f == "" {
		cwd, _ := os.Getwd()
		f = utils.DocValue(cwd)
	}
	cliLogger.Log.Infof("Start webserver on http://%s:%s, static file %s", a, p, f)
	err := http.ListenAndServe(string(a)+":"+string(p), http.FileServer(http.Dir(f)))
	if err != nil {
		cliLogger.Log.Error(err.Error())
	}
	return 0
}
