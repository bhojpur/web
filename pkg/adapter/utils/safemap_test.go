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

import "testing"

var safeMap *BhojpurMap

func TestNewBhojpurMap(t *testing.T) {
	safeMap = NewBhojpurMap()
	if safeMap == nil {
		t.Fatal("expected to return non-nil BhojpurMap", "got", safeMap)
	}
}

func TestSet(t *testing.T) {
	safeMap = NewBhojpurMap()
	if ok := safeMap.Set("bhojpur", 1); !ok {
		t.Error("expected", true, "got", false)
	}
}

func TestReSet(t *testing.T) {
	safeMap := NewBhojpurMap()
	if ok := safeMap.Set("bhojpur", 1); !ok {
		t.Error("expected", true, "got", false)
	}
	// set diff value
	if ok := safeMap.Set("bhojpur", -1); !ok {
		t.Error("expected", true, "got", false)
	}

	// set same value
	if ok := safeMap.Set("bhojpur", -1); ok {
		t.Error("expected", false, "got", true)
	}
}

func TestCheck(t *testing.T) {
	if exists := safeMap.Check("bhojpur"); !exists {
		t.Error("expected", true, "got", false)
	}
}

func TestGet(t *testing.T) {
	if val := safeMap.Get("bhojpur"); val.(int) != 1 {
		t.Error("expected value", 1, "got", val)
	}
}

func TestDelete(t *testing.T) {
	safeMap.Delete("bhojpur")
	if exists := safeMap.Check("bhojpur"); exists {
		t.Error("expected element to be deleted")
	}
}

func TestItems(t *testing.T) {
	safeMap := NewBhojpurMap()
	safeMap.Set("bhojpur", "hello")
	for k, v := range safeMap.Items() {
		key := k.(string)
		value := v.(string)
		if key != "bhojpur" {
			t.Error("expected the key should be bhojpur")
		}
		if value != "hello" {
			t.Error("expected the value should be hello")
		}
	}
}

func TestCount(t *testing.T) {
	if count := safeMap.Count(); count != 0 {
		t.Error("expected count to be", 0, "got", count)
	}
}
