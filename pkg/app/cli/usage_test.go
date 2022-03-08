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
	"bytes"
	"reflect"
	"testing"

	"github.com/bhojpur/web/pkg/app/errors"
)

func TestCommandUsage(t *testing.T) {
	cmd := &command{
		name: "foo bar",
		help: `
			Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
			eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
			ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
			aliquip ex ea commodo consequat. Duis aute irure dolor in
			reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
			pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
			culpa qui officia deserunt mollit anim id est laborum.
			`,
	}

	opts := []option{
		{
			name: "foo",
			help: `
				Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod
				tempor incididunt ut labore et dolore magna aliqua.
				`,
			envKey: "FOO",
			value:  reflect.ValueOf(42),
		},
		{
			name:   "bar",
			help:   "Bar option description.",
			envKey: "-",
			value:  reflect.ValueOf("bar"),
		},
		{
			name:   "alakazam",
			help:   "Alakazam option description.",
			envKey: "BAR",
			value:  reflect.ValueOf(0),
		},
	}

	w := bytes.NewBufferString("\n")
	usage := commandUsage(w, cmd, opts)
	usage()

	t.Log(w.String())
}

func TestCommandUsageIndex(t *testing.T) {
	m := commandManager{}

	m.register("foo", "bar").Help(`
	Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
	eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
	ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
	aliquip ex ea commodo consequat. Duis aute irure dolor in
	reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
	pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
	culpa qui officia deserunt mollit anim id est laborum.
	`)
	m.register("foo", "foo").Help("Foo lolilol.")
	m.register("foo", "buu").Help("A more simple help.")

	w := bytes.NewBufferString("\n")
	usage := commandUsageIndex(w, m.commands)
	usage()

	t.Log(w.String())
}

func TestPrintError(t *testing.T) {
	w := bytes.NewBufferString("\n")
	printError(w, errors.New("an error for testing printing"))
	t.Log(w.String())
}
