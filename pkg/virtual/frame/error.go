package frame

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
)

type ErrorType int

const (
	ErrorFrameSize ErrorType = iota
	ErrorProtocol
	ErrorProtocolStream
)

type Error struct {
	errorType ErrorType
	error
}

func (e Error) Type() ErrorType {
	return e.errorType
}

func (e Error) Err() error {
	return e.error
}

func frameSizeError(length uint32, frameName string) error {
	return &Error{ErrorFrameSize, fmt.Errorf("illegal %s frame length: 0x%x", frameName, length)}
}

func protoError(fmtstr string, args ...interface{}) error {
	return &Error{ErrorProtocol, fmt.Errorf(fmtstr, args...)}
}

func protoStreamError(fmtstr string, args ...interface{}) error {
	return &Error{ErrorProtocolStream, fmt.Errorf(fmtstr, args...)}
}
