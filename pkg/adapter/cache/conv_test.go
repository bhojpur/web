package cache

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

func TestGetString(t *testing.T) {
	var t1 = "test1"
	if "test1" != GetString(t1) {
		t.Error("get string from string error")
	}
	var t2 = []byte("test2")
	if "test2" != GetString(t2) {
		t.Error("get string from byte array error")
	}
	var t3 = 1
	if "1" != GetString(t3) {
		t.Error("get string from int error")
	}
	var t4 int64 = 1
	if "1" != GetString(t4) {
		t.Error("get string from int64 error")
	}
	var t5 = 1.1
	if "1.1" != GetString(t5) {
		t.Error("get string from float64 error")
	}

	if "" != GetString(nil) {
		t.Error("get string from nil error")
	}
}

func TestGetInt(t *testing.T) {
	var t1 = 1
	if 1 != GetInt(t1) {
		t.Error("get int from int error")
	}
	var t2 int32 = 32
	if 32 != GetInt(t2) {
		t.Error("get int from int32 error")
	}
	var t3 int64 = 64
	if 64 != GetInt(t3) {
		t.Error("get int from int64 error")
	}
	var t4 = "128"
	if 128 != GetInt(t4) {
		t.Error("get int from num string error")
	}
	if 0 != GetInt(nil) {
		t.Error("get int from nil error")
	}
}

func TestGetInt64(t *testing.T) {
	var i int64 = 1
	var t1 = 1
	if i != GetInt64(t1) {
		t.Error("get int64 from int error")
	}
	var t2 int32 = 1
	if i != GetInt64(t2) {
		t.Error("get int64 from int32 error")
	}
	var t3 int64 = 1
	if i != GetInt64(t3) {
		t.Error("get int64 from int64 error")
	}
	var t4 = "1"
	if i != GetInt64(t4) {
		t.Error("get int64 from num string error")
	}
	if 0 != GetInt64(nil) {
		t.Error("get int64 from nil")
	}
}

func TestGetFloat64(t *testing.T) {
	var f = 1.11
	var t1 float32 = 1.11
	if f != GetFloat64(t1) {
		t.Error("get float64 from float32 error")
	}
	var t2 = 1.11
	if f != GetFloat64(t2) {
		t.Error("get float64 from float64 error")
	}
	var t3 = "1.11"
	if f != GetFloat64(t3) {
		t.Error("get float64 from string error")
	}

	var f2 float64 = 1
	var t4 = 1
	if f2 != GetFloat64(t4) {
		t.Error("get float64 from int error")
	}

	if 0 != GetFloat64(nil) {
		t.Error("get float64 from nil error")
	}
}

func TestGetBool(t *testing.T) {
	var t1 = true
	if !GetBool(t1) {
		t.Error("get bool from bool error")
	}
	var t2 = "true"
	if !GetBool(t2) {
		t.Error("get bool from string error")
	}
	if GetBool(nil) {
		t.Error("get bool from nil error")
	}
}

func byteArrayEquals(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
