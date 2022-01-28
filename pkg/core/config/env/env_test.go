package env

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
	"os"
	"testing"
)

func TestEnvGet(t *testing.T) {
	gopath := Get("GOPATH", "")
	if gopath != os.Getenv("GOPATH") {
		t.Error("expected GOPATH not empty.")
	}

	noExistVar := Get("NOEXISTVAR", "foo")
	if noExistVar != "foo" {
		t.Errorf("expected NOEXISTVAR to equal foo, got %s.", noExistVar)
	}
}

func TestEnvMustGet(t *testing.T) {
	gopath, err := MustGet("GOPATH")
	if err != nil {
		t.Error(err)
	}

	if gopath != os.Getenv("GOPATH") {
		t.Errorf("expected GOPATH to be the same, got %s.", gopath)
	}

	_, err = MustGet("NOEXISTVAR")
	if err == nil {
		t.Error("expected error to be non-nil")
	}
}

func TestEnvSet(t *testing.T) {
	Set("MYVAR", "foo")
	myVar := Get("MYVAR", "bar")
	if myVar != "foo" {
		t.Errorf("expected MYVAR to equal foo, got %s.", myVar)
	}
}

func TestEnvMustSet(t *testing.T) {
	err := MustSet("FOO", "bar")
	if err != nil {
		t.Error(err)
	}

	fooVar := os.Getenv("FOO")
	if fooVar != "bar" {
		t.Errorf("expected FOO variable to equal bar, got %s.", fooVar)
	}
}

func TestEnvGetAll(t *testing.T) {
	envMap := GetAll()
	if len(envMap) == 0 {
		t.Error("expected environment not empty.")
	}
}
