package update

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
	"flag"
	"os"
	"os/exec"

	"github.com/bhojpur/web/cmd/utility/commands"
	"github.com/bhojpur/web/pkg/client/config"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/utils"
)

var CmdUpdate = &commands.Command{
	UsageLine: "update",
	Short:     "Update Bhojpur Web CLI",
	Long: `
Automatic run command "go get -u github.com/bhojpur/web" for selfupdate
`,
	Run: updateBhojpurWeb,
}

func init() {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	CmdUpdate.Flag = *fs
	commands.AvailableCommands = append(commands.AvailableCommands, CmdUpdate)
}

func updateBhojpurWeb(cmd *commands.Command, args []string) int {
	cliLogger.Log.Info("Updating")
	bhojpurWebPath := config.GitRemotePath
	cmdUp := exec.Command("go", "get", "-u", bhojpurWebPath)
	cmdUp.Stdout = os.Stdout
	cmdUp.Stderr = os.Stderr
	if err := cmdUp.Run(); err != nil {
		cliLogger.Log.Warnf("Run cmd err:%s", err)
	}
	// update the Time when updateBhojpurWeb every time
	utils.UpdateLastPublishedTime()
	return 0
}
