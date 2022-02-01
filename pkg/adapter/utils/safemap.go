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

// BeeMap is a map with lock
type BeeMap utils.BeeMap

// NewBeeMap return new safemap
func NewBeeMap() *BeeMap {
	return (*BeeMap)(utils.NewBeeMap())
}

// Get from maps return the k's value
func (m *BeeMap) Get(k interface{}) interface{} {
	return (*utils.BeeMap)(m).Get(k)
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *BeeMap) Set(k interface{}, v interface{}) bool {
	return (*utils.BeeMap)(m).Set(k, v)
}

// Check Returns true if k is exist in the map.
func (m *BeeMap) Check(k interface{}) bool {
	return (*utils.BeeMap)(m).Check(k)
}

// Delete the given key and value.
func (m *BeeMap) Delete(k interface{}) {
	(*utils.BeeMap)(m).Delete(k)
}

// Items returns all items in safemap.
func (m *BeeMap) Items() map[interface{}]interface{} {
	return (*utils.BeeMap)(m).Items()
}

// Count returns the number of items within the map.
func (m *BeeMap) Count() int {
	return (*utils.BeeMap)(m).Count()
}
