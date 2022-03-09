package run

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
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/bhojpur/web/pkg/client/config"
	cliLogger "github.com/bhojpur/web/pkg/client/logger"
	"github.com/bhojpur/web/pkg/client/logger/colors"
	"github.com/bhojpur/web/pkg/client/utils"
	"github.com/fsnotify/fsnotify"
)

var (
	cmd                 *exec.Cmd
	state               sync.Mutex
	eventTime           = make(map[string]int64)
	scheduleTime        time.Time
	watchExts           = config.Conf.WatchExts
	watchExtsStatic     = config.Conf.WatchExtsStatic
	ignoredFilesRegExps = []string{
		`.#(\w+).go$`,
		`.(\w+).go.swp$`,
		`(\w+).go~$`,
		`(\w+).tmp$`,
		`commentsRouter_controllers.go$`,
	}
)

// NewWatcher starts an fsnotify Watcher on the specified paths
func NewWatcher(paths []string, files []string, isgenerate bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		cliLogger.Log.Fatalf("Failed to create watcher: %s", err)
	}

	go func() {
		for {
			select {
			case e := <-watcher.Events:
				isBuild := true

				if ifStaticFile(e.Name) && config.Conf.EnableReload {
					sendReload(e.String())
					continue
				}
				// Skip ignored files
				if shouldIgnoreFile(e.Name) {
					continue
				}
				if !shouldWatchFileWithExtension(e.Name) {
					continue
				}

				mt := utils.GetFileModTime(e.Name)
				if t := eventTime[e.Name]; mt == t {
					cliLogger.Log.Hintf(colors.Bold("Skipping: ")+"%s", e.String())
					isBuild = false
				}

				eventTime[e.Name] = mt

				if isBuild {
					cliLogger.Log.Hintf("Event fired: %s", e)
					go func() {
						// Wait 1s before autobuild until there is no file change.
						scheduleTime = time.Now().Add(1 * time.Second)
						time.Sleep(time.Until(scheduleTime))
						AutoBuild(files, isgenerate)

						if config.Conf.EnableReload {
							// Wait 100ms more before refreshing the browser
							time.Sleep(100 * time.Millisecond)
							sendReload(e.String())
						}
					}()
				}
			case err := <-watcher.Errors:
				cliLogger.Log.Warnf("Watcher error: %s", err.Error()) // No need to exit here
			}
		}
	}()

	cliLogger.Log.Info("Initializing watcher...")
	for _, path := range paths {
		cliLogger.Log.Hintf(colors.Bold("Watching: ")+"%s", path)
		err = watcher.Add(path)
		if err != nil {
			cliLogger.Log.Fatalf("Failed to watch directory: %s", err)
		}
	}
}

// AutoBuild builds the specified set of files
func AutoBuild(files []string, isgenerate bool) {
	state.Lock()
	defer state.Unlock()

	os.Chdir(currpath)

	cmdName := "go"

	var (
		err    error
		stderr bytes.Buffer
	)
	// For applications use full import path like "github.com/.../.."
	// are able to use "go install" to reduce build time.
	if config.Conf.GoInstall {
		icmd := exec.Command(cmdName, "install", "-v")
		icmd.Stdout = os.Stdout
		icmd.Stderr = os.Stderr
		icmd.Env = append(os.Environ(), "GOGC=off")
		icmd.Run()
	}

	if isgenerate {
		cliLogger.Log.Info("Generating the docs...")
		icmd := exec.Command("webutl", "generate", "docs")
		icmd.Env = append(os.Environ(), "GOGC=off")
		err = icmd.Run()
		if err != nil {
			utils.Notify("", "Failed to generate the docs.")
			cliLogger.Log.Errorf("Failed to generate the docs.")
			return
		}
		cliLogger.Log.Success("Docs generated!")
	}
	appName := appname
	if err == nil {

		if runtime.GOOS == "windows" {
			appName += ".exe"
		}

		args := []string{"build"}
		args = append(args, "-o", appName)
		if buildTags != "" {
			args = append(args, "-tags", buildTags)
		}
		if buildLDFlags != "" {
			args = append(args, "-ldflags", buildLDFlags)
		}
		args = append(args, files...)

		bcmd := exec.Command(cmdName, args...)
		bcmd.Env = append(os.Environ(), "GOGC=off")
		bcmd.Stderr = &stderr
		err = bcmd.Run()
		if err != nil {
			utils.Notify(stderr.String(), "Build Failed")
			cliLogger.Log.Errorf("Failed to build the application: %s", stderr.String())
			return
		}
	}

	cliLogger.Log.Success("Built Successfully!")
	Restart(appName)
}

// Kill kills the running command process
func Kill() {
	defer func() {
		if e := recover(); e != nil {
			cliLogger.Log.Infof("Kill recover: %s", e)
		}
	}()
	if cmd != nil && cmd.Process != nil {
		// Windows does not support Interrupt
		if runtime.GOOS == "windows" {
			cmd.Process.Signal(os.Kill)
		} else {
			cmd.Process.Signal(os.Interrupt)
		}

		ch := make(chan struct{}, 1)
		go func() {
			cmd.Wait()
			ch <- struct{}{}
		}()

		select {
		case <-ch:
			return
		case <-time.After(10 * time.Second):
			cliLogger.Log.Info("Timeout. Force kill cmd process")
			err := cmd.Process.Kill()
			if err != nil {
				cliLogger.Log.Errorf("Error while killing cmd process: %s", err)
			}
			return
		}
	}
}

// Restart kills the running command process and starts it again
func Restart(appname string) {
	cliLogger.Log.Debugf("Kill running process", utils.FILE(), utils.LINE())
	Kill()
	go Start(appname)
}

// Start starts the command process
func Start(appname string) {
	cliLogger.Log.Infof("Restarting '%s'...", appname)
	if !strings.Contains(appname, "./") {
		appname = "./" + appname
	}

	cmd = exec.Command(appname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runargs != "" {
		r := regexp.MustCompile("'.+'|\".+\"|\\S+")
		m := r.FindAllString(runargs, -1)
		cmd.Args = append([]string{appname}, m...)
	} else {
		cmd.Args = append([]string{appname}, config.Conf.CmdArgs...)
	}
	cmd.Env = append(os.Environ(), config.Conf.Envs...)

	go cmd.Run()
	cliLogger.Log.Successf("'%s' is running...", appname)
	started <- true
}

func ifStaticFile(filename string) bool {
	for _, s := range watchExtsStatic {
		if strings.HasSuffix(filename, s) {
			return true
		}
	}
	return false
}

// shouldIgnoreFile ignores filenames generated by Emacs, Vim or SublimeText.
// It returns true if the file should be ignored, false otherwise.
func shouldIgnoreFile(filename string) bool {
	for _, regex := range ignoredFilesRegExps {
		r, err := regexp.Compile(regex)
		if err != nil {
			cliLogger.Log.Fatalf("Could not compile regular expression: %s", err)
		}
		if r.MatchString(filename) {
			return true
		}
		continue
	}
	return false
}

// shouldWatchFileWithExtension returns true if the name of the file
// hash a suffix that should be watched.
func shouldWatchFileWithExtension(name string) bool {
	for _, s := range watchExts {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}
