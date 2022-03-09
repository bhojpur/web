package rs

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
	"os/exec"
	"runtime"
	"time"

	"strings"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/cmd/utility/commands/version"
	"github.com/bhojpur/web/pkg/client/config"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/logger/colors"
	"github.com/bhojpur/web/pkg/client/utils"
)

var cmdRs = &commands.Command{
	UsageLine: "rs",
	Short:     "Run customized scripts",
	Long: `Run script allows you to run arbitrary commands using Bhojpur Web.
  Custom commands are provided from the "scripts" object inside bhojpur.json or Bhojpurfile.

  To run a custom command, use: {{"$ webutl rs mycmd ARGS" | bold}}
  {{if len .}}
{{"AVAILABLE SCRIPTS"|headline}}{{range $cmdName, $cmd := .}}
  {{$cmdName | bold}}
      {{$cmd}}{{end}}{{end}}
`,
	PreRun: func(cmd *commands.Command, args []string) { version.ShowShortVersionBanner() },
	Run:    runScript,
}

func init() {
	config.LoadConfig()
	cmdRs.Long = utils.TmplToString(cmdRs.Long, config.Conf.Scripts)
	commands.AvailableCommands = append(commands.AvailableCommands, cmdRs)
}

func runScript(cmd *commands.Command, args []string) int {
	if len(args) == 0 {
		cmd.Usage()
	}

	start := time.Now()
	script, args := args[0], args[1:]

	if c, exist := config.Conf.Scripts[script]; exist {
		command := customCommand{
			Name:    script,
			Command: c,
			Args:    args,
		}
		if err := command.run(); err != nil {
			cliLogger.Log.Error(err.Error())
		}
	} else {
		cliLogger.Log.Errorf("Command '%s' not found in Bhojpurfile/bhojpur.json", script)
	}
	elapsed := time.Since(start)
	fmt.Println(colors.GreenBold(fmt.Sprintf("Finished in %s.", elapsed)))
	return 0
}

type customCommand struct {
	Name    string
	Command string
	Args    []string
}

func (c *customCommand) run() error {
	cliLogger.Log.Info(colors.GreenBold(fmt.Sprintf("Running '%s'...", c.Name)))
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin", "linux":
		args := append([]string{c.Command}, c.Args...)
		cmd = exec.Command("sh", "-c", strings.Join(args, " "))
	case "windows":
		args := append([]string{c.Command}, c.Args...)
		cmd = exec.Command("cmd", "/C", strings.Join(args, " "))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
