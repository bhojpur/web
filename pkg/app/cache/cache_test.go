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
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func testCache(t *testing.T, c Cache) {
	ctx := context.TODO()

	_, isCached := c.Get(ctx, "/foo")
	require.False(t, isCached)
	require.Zero(t, c.Len())
	require.Zero(t, c.Size())

	c.Set(ctx, "/foo", Bytes("foo"))
	require.Equal(t, 1, c.Len())
	require.Equal(t, 3, c.Size())

	c.Set(ctx, "/bar", String("bar"))
	require.Equal(t, 2, c.Len())
	require.Equal(t, 6, c.Size())

	foo, isCached := c.Get(ctx, "/foo")
	require.True(t, isCached)
	require.Equal(t, Bytes("foo"), foo.(Bytes))

	bar, isCached := c.Get(ctx, "/bar")
	require.True(t, isCached)
	require.Equal(t, String("bar"), bar.(String))

	c.Del(ctx, "/foo")
	require.Equal(t, 1, c.Len())
	require.Equal(t, 3, c.Size())

	foo, isCached = c.Get(ctx, "/foo")
	require.False(t, isCached)
	require.Nil(t, foo)
}

func TestItemSize(t *testing.T) {
	utests := []struct {
		scenario     string
		item         Item
		expectedSize int
	}{
		{
			scenario:     "bytes",
			item:         Bytes("boo"),
			expectedSize: 3,
		},
		{
			scenario:     "string",
			item:         String("hello"),
			expectedSize: 5,
		},
		{
			scenario:     "int",
			item:         Int(42),
			expectedSize: int(unsafe.Sizeof(42)),
		},
		{
			scenario:     "float",
			item:         Float(42.0),
			expectedSize: int(unsafe.Sizeof(42.1)),
		},
	}

	for _, u := range utests {
		t.Run(u.scenario, func(t *testing.T) {
			require.Equal(t, u.expectedSize, u.item.Size())
		})
	}
}
