package utils

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
)

func TestCompareGoVersion(t *testing.T) {
	targetVersion := "go1.8"
	if compareGoVersion("go1.12.4", targetVersion) != 1 {
		t.Error("should be 1")
	}

	if compareGoVersion("go1.8.7", targetVersion) != 1 {
		t.Error("should be 1")
	}

	if compareGoVersion("go1.8", targetVersion) != 0 {
		t.Error("should be 0")
	}

	if compareGoVersion("go1.7.6", targetVersion) != -1 {
		t.Error("should be -1")
	}

	if compareGoVersion("go1.12.1rc1", targetVersion) != 1 {
		t.Error("should be 1")
	}

	if compareGoVersion("go1.8rc1", targetVersion) != 0 {
		t.Error("should be 0")
	}

	if compareGoVersion("go1.7rc1", targetVersion) != -1 {
		t.Error("should be -1")
	}
}
