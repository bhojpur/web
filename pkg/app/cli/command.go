package cli

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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/web/pkg/app/errors"
)

var (
	errNoRootCmd = errors.New("no root command")
)

type command struct {
	help    string
	name    string
	options interface{}
}

func (c *command) Help(h string) Command {
	c.help = h
	return c
}

func (c *command) Options(o interface{}) Command {
	c.options = o
	return c
}

type commandManager struct {
	out      io.Writer
	commands map[string]Command
}

func (m *commandManager) register(cmd ...string) Command {
	k := commandString(cmd...)
	c := &command{name: k}

	if m.commands == nil {
		m.commands = make(map[string]Command)
	}

	m.commands[k] = c
	return c
}

func (m *commandManager) parse(args ...string) (string, func(), error) {
	cmdslice, optsSlice := splitCommand(args)
	k := commandString(cmdslice...)

	i, ok := m.commands[k]
	if !ok {
		err := errNoRootCmd
		if k != "" {
			err = errors.New("unknown command").Tag("command-name", k)
		}
		return "", commandUsageIndex(m.out, m.commands), err
	}
	cmd := i.(*command)

	programName := filepath.Base(os.Args[0])
	flags := flag.NewFlagSet(commandString(programName, k), flag.ContinueOnError)
	flags.SetOutput(writerNoop{})

	optsParser := optionParser{flags: flags}
	opts, err := optsParser.parse(cmd.options)
	if err != nil {
		return "", nil, errors.New("parsing options failed").
			Tag("command-name", k).
			Wrap(err)
	}

	usage := commandUsage(m.out, cmd, opts)
	flags.Usage = func() {}
	return cmd.name, usage, flags.Parse(optsSlice)
}

func commandString(cmd ...string) string {
	clean := make([]string, 0, len(cmd))
	for _, c := range cmd {
		if c == "" {
			continue
		}
		clean = append(clean, strings.TrimSpace(c))
	}
	return strings.Join(clean, " ")
}

func splitCommand(args []string) (cmd, opts []string) {
	end := commandEndIndex(args)
	return args[:end], args[end:]
}

func commandEndIndex(args []string) int {
	i := 0
	for i < len(args) {
		if strings.HasPrefix(args[i], "-") {
			return i
		}
		i++
	}
	return i
}

type writerNoop struct{}

func (w writerNoop) Write([]byte) (int, error) {
	return 0, nil
}
