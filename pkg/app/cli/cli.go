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

// It is a package to build CLI programs. For example:
//  func main() {
//      opts := struct {
// 	        Foo string `env:"FOO" help:"A simple string."`
// 		    Bar int    `env:"BAR" help:"A simple integer."`
//      }{
// 		    Foo: "foo",
// 		    Bar: 42,
//      }
//
//      ctx, cancel := cli.ContextWithSignals(context.Background(),
//          os.Interrupt,
// 		    syscall.SIGTERM,
//      )
//      defer cancel()
//
//      cli.Register().
// 		    Help("A simple command").
// 		    Options(&opts)
//
//      cli.Register("hello").
// 		    Help("A sub command").
// 		    Options(&opts)
//
//      cli.Register("world").
// 		    Help("Another sub command").
// 		    Options(&opts)
//
//      switch cli.Load() {
// 	    case "hello":
// 		    helloCmd(ctx, opts)
//
//      case "world":
// 		    worldCmd(ctx, opts)
//
// 	    default:
// 		    defaultCmd(ctx, opts)
// 	    }
//  }
//

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
)

var (
	defaultManager = commandManager{out: os.Stderr}
	currentUsage   func()
	exitOnError    = true
	programArgs    = os.Args[1:]
)

// Command is the interface that describes a command.
type Command interface {
	// Sets the command help description.
	Help(string) Command

	// Sets the command options with the given receiver. The receiver must be a
	// pointer to a struct.
	Options(interface{}) Command
}

// Register registers and returns the named command.
func Register(cmd ...string) Command {
	return defaultManager.register(cmd...)
}

// Load loads the registered command that matches the program args. If defined,
// environment variables and flags are loaded in the command options.
//
// It prints the command usage and exits the program with code -1 when an error
// occurs.
func Load() (cmd string) {
	cmd, usage, err := defaultManager.parse(programArgs...)
	currentUsage = usage

	if err != nil {
		if errors.Is(err, errNoRootCmd) && err != flag.ErrHelp {
			printError(defaultManager.out, err)
		}

		if usage != nil {
			usage()
		}

		if exitOnError {
			os.Exit(-1)
		}
		panic(err)
	}

	return cmd
}

// Usage prints the loaded command usage. It panics when called before the Load
// function.
func Usage() {
	if currentUsage == nil {
		panic("usage func is called before load func")
	}
	currentUsage()
}

// Error prints the given error and exit the program with code -1.
func Error(err error) {
	printError(defaultManager.out, err)
	if exitOnError {
		os.Exit(-1)
	}
}

// ContextWithSignals returns a copy of the parent context that gets canceled
// when one of the specified signals is emitted. If no signals are provided, all
// incoming signals will cancel the context.
//
// Canceling this context releases resources associated with it, so code should
// call cancel as soon as the operations running in this Context complete.
func ContextWithSignals(parent context.Context, sig ...os.Signal) (ctx context.Context, cancel func()) {
	ctx, cancel = context.WithCancel(parent)
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)

	go func() {
		<-c
		close(c)
		cancel()
	}()

	return ctx, cancel
}
