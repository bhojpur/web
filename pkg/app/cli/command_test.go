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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommandManagerNoCommands(t *testing.T) {
	m := commandManager{}
	_, _, err := m.parse()
	require.Error(t, err)
	t.Log("error:", err)
}

func TestCommandManagerRegisterRoot(t *testing.T) {
	opts := struct {
		Int    int
		String string
	}{
		Int:    42,
		String: "foo",
	}

	m := commandManager{}
	m.register().
		Help("Test command.").
		Options(&opts)

	cmd, _, err := m.parse()
	require.NoError(t, err)
	require.Empty(t, cmd)
	require.Equal(t, 42, opts.Int)
	require.Equal(t, "foo", opts.String)
}

func TestCommandManagerRegisterRootWithOptions(t *testing.T) {
	opts := struct {
		Int    int
		String string
	}{}

	m := commandManager{}
	m.register().
		Help("Test command.").
		Options(&opts)

	cmd, _, err := m.parse("-int", "21", "-string", "bar")
	require.NoError(t, err)
	require.Empty(t, cmd)
	require.Equal(t, 21, opts.Int)
	require.Equal(t, "bar", opts.String)
}

func TestCommandManagerRegisterWithBadOptions(t *testing.T) {
	opts := struct {
		Int    int
		String string
	}{}

	m := commandManager{}
	m.register().
		Help("Test command.").
		Options(opts)

	_, _, err := m.parse("-int", "21", "-string", "bar")
	require.Error(t, err)
	t.Log("error:", err)
}

func TestCommandManagerRegisterMultiple(t *testing.T) {
	optsBar := struct {
		Int    int
		String string
	}{}

	optsBoo := optsBar

	m := commandManager{}

	m.register("foo", "bar").
		Help("Test command.").
		Options(&optsBar)

	m.register("foo", "boo").
		Help("Test command 2.").
		Options(&optsBoo)

	cmd, _, err := m.parse("foo", "boo", "-int", "21", "-string", "bar")
	require.NoError(t, err)
	require.Equal(t, "foo boo", cmd)
	require.Equal(t, 21, optsBoo.Int)
	require.Equal(t, "bar", optsBoo.String)
	require.Zero(t, optsBar.Int)
	require.Zero(t, optsBar.String)
}

func TestCommandeString(t *testing.T) {
	tests := []struct {
		scenario string
		command  []string
		expected string
	}{
		{
			scenario: "empty command",
		},
		{
			scenario: "command with one element",
			command:  []string{"test"},
			expected: "test",
		},
		{
			scenario: "command with multiple elements",
			command:  []string{"test", "foo", "bar"},
			expected: "test foo bar",
		},
		{
			scenario: "command with elements with trailing spaces",
			command:  []string{"\ttest    ", "foo\n", "     bar"},
			expected: "test foo bar",
		},
		{
			scenario: "command with elements with empty elements",
			command:  []string{"", "foo", "bar"},
			expected: "foo bar",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			cmd := commandString(test.command...)
			require.Equal(t, test.expected, cmd)
		})
	}
}

func TestCommandEndIndex(t *testing.T) {
	tests := []struct {
		scenario      string
		args          []string
		expectedIndex int
	}{
		{
			scenario:      "args with an option at the beginning",
			args:          []string{"-v", "foo", "bar"},
			expectedIndex: 0,
		},
		{
			scenario:      "args with an option at the end",
			args:          []string{"foo", "bar", "-v"},
			expectedIndex: 2,
		},
		{
			scenario:      "args with an option at the middle",
			args:          []string{"foo", "foo-bar", "-v", "bar"},
			expectedIndex: 2,
		},
		{
			scenario:      "args without option",
			args:          []string{"foo", "foo-bar", "bar"},
			expectedIndex: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			idx := commandEndIndex(test.args)
			require.Equal(t, test.expectedIndex, idx)
		})
	}
}
