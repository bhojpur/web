package berror

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
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// code, msg
const errFmt = "ERROR-%d, %s"

// Err returns an error representing c and msg. If c is OK, returns nil.
func Error(c Code, msg string) error {
	return fmt.Errorf(errFmt, c.Code(), msg)
}

// Errorf returns error
func Errorf(c Code, format string, a ...interface{}) error {
	return Error(c, fmt.Sprintf(format, a...))
}

func Wrap(err error, c Code, msg string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, fmt.Sprintf(errFmt, c.Code(), msg))
}

func Wrapf(err error, c Code, format string, a ...interface{}) error {
	return Wrap(err, c, fmt.Sprintf(format, a...))
}

// FromError is very simple. It just parse error msg and check whether code has been register
// if code not being register, return unknown
// if err.Error() is not valid bhojpur error code, return unknown
func FromError(err error) (Code, bool) {
	msg := err.Error()
	codeSeg := strings.SplitN(msg, ",", 2)
	if strings.HasPrefix(codeSeg[0], "ERROR-") {
		codeStr := strings.SplitN(codeSeg[0], "-", 2)
		if len(codeStr) < 2 {
			return Unknown, false
		}
		codeInt, e := strconv.ParseUint(codeStr[1], 10, 32)
		if e != nil {
			return Unknown, false
		}
		if code, ok := defaultCodeRegistry.Get(uint32(codeInt)); ok {
			return code, true
		}
	}
	return Unknown, false
}
