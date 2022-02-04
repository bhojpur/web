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
	"github.com/bhojpur/web/pkg/core/utils"
)

// BhojpurMap is a map with lock
type BhojpurMap utils.BhojpurMap

// NewBhojpurMap return new safemap
func NewBhojpurMap() *BhojpurMap {
	return (*BhojpurMap)(utils.NewBhojpurMap())
}

// Get from maps return the k's value
func (m *BhojpurMap) Get(k interface{}) interface{} {
	return (*utils.BhojpurMap)(m).Get(k)
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BhojpurMap) Set(k interface{}, v interface{}) bool {
	return (*utils.BhojpurMap)(m).Set(k, v)
}

// Check Returns true if k is exist in the map.
func (m *BhojpurMap) Check(k interface{}) bool {
	return (*utils.BhojpurMap)(m).Check(k)
}

// Delete the given key and value.
func (m *BhojpurMap) Delete(k interface{}) {
	(*utils.BhojpurMap)(m).Delete(k)
}

// Items returns all items in safemap.
func (m *BhojpurMap) Items() map[interface{}]interface{} {
	return (*utils.BhojpurMap)(m).Items()
}

// Count returns the number of items within the map.
func (m *BhojpurMap) Count() int {
	return (*utils.BhojpurMap)(m).Count()
}
