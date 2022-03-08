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
	"context"
	"reflect"
)

// Cache is the interface that describes a cache.
type Cache interface {
	// Get returns the item with the given key, otherwise returns false.
	Get(ctx context.Context, key string) (Item, bool)

	// Set sets the item at the given key.
	Set(ctx context.Context, key string, i Item)

	// Deletes the item at the given key.
	Del(ctx context.Context, key string)

	// The number of items in the cache.
	Len() int

	// The size in bytes.
	Size() int
}

// Item is the interface that describes a cacheable item.
type Item interface {
	// The size that the item occupies in a cache.
	Size() int
}

// Bytes represents a cacheable byte slice.
type Bytes []byte

func (b Bytes) Size() int {
	return len(b)
}

type String string

func (s String) Size() int {
	return len(s)
}

type Int int

func (i Int) Size() int {
	return intSize
}

type Float float64

func (f Float) Size() int {
	return floatSize
}

var (
	intSize   = int(reflect.TypeOf(42).Size())
	floatSize = int(reflect.TypeOf(23.42).Size())
)
