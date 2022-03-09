package utils

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
	"os/exec"
	"strconv"
	"strings"

	"runtime"

	"github.com/bhojpur/web/pkg/client/config"
)

const appName = "webutl"

func Notify(text, title string) {
	if !config.Conf.EnableNotification {
		return
	}
	switch runtime.GOOS {
	case "darwin":
		osxNotify(text, title)
	case "linux":
		linuxNotify(text, title)
	case "windows":
		windowsNotify(text, title)
	}
}

func osxNotify(text, title string) {
	var cmd *exec.Cmd
	if existTerminalNotifier() {
		cmd = exec.Command("terminal-notifier", "-title", appName, "-message", text, "-subtitle", title)
	} else if MacOSVersionSupport() {
		notification := fmt.Sprintf("display notification \"%s\" with title \"%s\" subtitle \"%s\"", text, appName, title)
		cmd = exec.Command("osascript", "-e", notification)
	} else {
		cmd = exec.Command("growlnotify", "-n", appName, "-m", title)
	}
	cmd.Run()
}

func windowsNotify(text, title string) {
	exec.Command("growlnotify", "/i:", "", "/t:", title, text).Run()
}

func linuxNotify(text, title string) {
	exec.Command("notify-send", "-i", "", title, text).Run()
}

func existTerminalNotifier() bool {
	cmd := exec.Command("which", "terminal-notifier")
	err := cmd.Start()
	if err != nil {
		return false
	}
	err = cmd.Wait()
	return err != nil
}

func MacOSVersionSupport() bool {
	cmd := exec.Command("sw_vers", "-productVersion")
	check, _ := cmd.Output()
	version := strings.Split(string(check), ".")
	major, _ := strconv.Atoi(version[0])
	minor, _ := strconv.Atoi(version[1])
	if major < 10 || (major == 10 && minor < 9) {
		return false
	}
	return true
}
