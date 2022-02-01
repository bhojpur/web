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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCode1 = DefineCode(1, "unit_test", "TestError", "Hello, test code1")

var testErr = errors.New("hello, this is error")

func TestErrorf(t *testing.T) {
	msg := Errorf(testCode1, "errorf %s", "aaaa")
	assert.NotNil(t, msg)
	assert.Equal(t, "ERROR-1, errorf aaaa", msg.Error())
}

func TestWrapf(t *testing.T) {
	err := Wrapf(testErr, testCode1, "Wrapf %s", "aaaa")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(err, testErr))
}

func TestFromError(t *testing.T) {
	err := errors.New("ERROR-1, errorf aaaa")
	code, ok := FromError(err)
	assert.True(t, ok)
	assert.Equal(t, testCode1, code)
	assert.Equal(t, "unit_test", code.Module())
	assert.Equal(t, "Hello, test code1", code.Desc())

	err = errors.New("not bhojpur error")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR-2, not register")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR-aaa, invalid code")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("aaaaaaaaaaaaaa")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR-2-3, invalid error")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)

	err = errors.New("ERROR, invalid error")
	code, ok = FromError(err)
	assert.False(t, ok)
	assert.Equal(t, Unknown, code)
}
