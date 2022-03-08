package app

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
	"runtime"
	"strings"
)

var (
	// DefaultLogger is the logger used to log info and errors.
	DefaultLogger func(format string, v ...interface{})

	defaultColor string
	errorColor   string
	infoColor    string
)

func init() {
	goarch := runtime.GOARCH
	if goarch == "wasm" {
		DefaultLogger = clientLog
		return
	}

	if goarch != "window" {
		defaultColor = "\033[00m"
		errorColor = "\033[91m"
		infoColor = "\033[94m"
	}
	DefaultLogger = serverLog
}

// Log logs using the default formats for its operands. Spaces are always added
// between operands.
func Log(v ...interface{}) {
	var b strings.Builder
	for i := 0; i < len(v); i++ {
		if i != 0 {
			b.WriteByte(' ')
		}
		b.WriteString("%v")
	}
	Logf(b.String(), v...)
}

// Logf logs according to a format specifier.
func Logf(format string, v ...interface{}) {
	DefaultLogger(format, v...)
}

func serverLog(format string, v ...interface{}) {
	errorLevel := false

	for _, a := range v {
		if _, ok := a.(error); ok {
			errorLevel = true
			break
		}
	}

	if errorLevel {
		fmt.Printf(errorColor+"ERROR ‣ "+defaultColor+format+"\n", v...)
		return
	}

	fmt.Printf(infoColor+"INFO ‣ "+defaultColor+format+"\n", v...)
}

func clientLog(format string, v ...interface{}) {
	isErrorLevel := false
	for _, a := range v {
		if _, isErr := a.(error); isErr {
			isErrorLevel = true
			break
		}
	}

	if isErrorLevel {
		Window().Get("console").Call("error", fmt.Sprintf(format, v...))
		return
	}
	Window().Get("console").Call("log", fmt.Sprintf(format, v...))
}
