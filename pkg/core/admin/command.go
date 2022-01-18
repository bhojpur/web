package admin

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
	"github.com/pkg/errors"
)

// Command is an experimental interface
// We try to use this to decouple modules
// All other modules depends on this, and they register the command they support
// We may change the API in the future, so be careful about this.
type Command interface {
	Execute(params ...interface{}) *Result
}

var CommandNotFound = errors.New("Command not found")

type Result struct {
	// Status is the same as http.Status
	Status  int
	Error   error
	Content interface{}
}

func (r *Result) IsSuccess() bool {
	return r.Status >= 200 && r.Status < 300
}

// CommandRegistry stores all commands
// name => command
type moduleCommands map[string]Command

// Get returns command with the name
func (m moduleCommands) Get(name string) Command {
	c, ok := m[name]
	if ok {
		return c
	}
	return &doNothingCommand{}
}

// module name => moduleCommand
type commandRegistry map[string]moduleCommands

// Get returns module's commands
func (c commandRegistry) Get(moduleName string) moduleCommands {
	if mcs, ok := c[moduleName]; ok {
		return mcs
	}
	res := make(moduleCommands)
	c[moduleName] = res
	return res
}

var cmdRegistry = make(commandRegistry)

// RegisterCommand is not thread-safe
// do not use it in concurrent case
func RegisterCommand(module string, commandName string, command Command) {
	cmdRegistry.Get(module)[commandName] = command
}

func GetCommand(module string, cmdName string) Command {
	return cmdRegistry.Get(module).Get(cmdName)
}

type doNothingCommand struct{}

func (d *doNothingCommand) Execute(params ...interface{}) *Result {
	return &Result{
		Status: 404,
		Error:  CommandNotFound,
	}
}
