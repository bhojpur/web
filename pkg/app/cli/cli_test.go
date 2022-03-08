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
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/bhojpur/web/pkg/app/errors"
	"github.com/stretchr/testify/require"
)

func init() {
	exitOnError = false
}

func TestCliSuccess(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	defaultManager.commands = nil
	programArgs = nil

	Register().Help("A test command")
	cmd := Load()
	require.Empty(t, cmd)

	Usage()
	t.Log(w.String())
}

func TestCliIndex(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	defaultManager.commands = nil
	programArgs = []string{"foo", "test"}

	Register("foo", "bar").Help("A test command")
	Register("foo", "buu").Help("Another test command")

	defer func() {
		recover()
		t.Log(w.String())
	}()

	Load()
	t.Fail()
}

func TestCliCmdBadOption(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	defaultManager.commands = nil
	programArgs = []string{"-duration", "[x_x]"}

	opts := struct {
		Duration time.Duration
	}{}

	Register().Options(&opts)

	defer func() {
		recover()
		t.Log(w.String())
	}()

	Load()
	t.Fail()
}

func TestUsagePanic(t *testing.T) {
	currentUsage = nil
	require.Panics(t, func() {
		Usage()
	})
}

func TestError(t *testing.T) {
	w := bytes.NewBufferString("\n")
	defaultManager.out = w
	Error(errors.New("error error critical error"))
	t.Log(w.String())
}

func TestContextWithSignals(t *testing.T) {
	ctx, cancel := ContextWithSignals(context.TODO(), os.Interrupt)
	defer cancel()

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	<-ctx.Done()
	require.Error(t, ctx.Err())
}
