package app

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
	"time"

	"github.com/stretchr/testify/require"
)

func TestPreRenderLRUCache(t *testing.T) {
	testPreRenderCache(t, NewPreRenderLRUCache(100, time.Second))
}

func TestPreRenderCache(t *testing.T) {
	testPreRenderCache(t, newPreRenderCache(1))
}

func testPreRenderCache(t *testing.T, c PreRenderCache) {
	ctx := context.TODO()

	i := PreRenderedItem{
		Path:            "/test",
		ContentType:     "text/html",
		ContentEncoding: "gzip",
		Body:            []byte("test"),
	}

	ic, ok := c.Get(ctx, i.Path)
	require.Zero(t, ic)
	require.False(t, ok)

	c.Set(ctx, i)
	ic, ok = c.Get(ctx, i.Path)
	require.True(t, ok)
	require.Equal(t, i, ic)
}

func TestPreRenderLRUCacheExpire(t *testing.T) {
	ctx := context.TODO()
	evictCalled := false
	onEvict := func(string, PreRenderedItem) { evictCalled = true }

	c := NewPreRenderLRUCache(16, -time.Second, onEvict).(*preRenderLRUCache)

	items := []PreRenderedItem{
		{
			Path: "/test1",
			Body: []byte("test"),
		},
		{
			Path: "/test2",
			Body: []byte("test"),
		},
		{
			Path: "/test3",
			Body: []byte("test"),
		},
		{
			Path: "/test4",
			Body: []byte("test"),
		},
	}

	for _, i := range items {
		c.Set(ctx, i)
	}
	require.Equal(t, 4, c.Len())
	require.Equal(t, 16, c.Size())

	for _, i := range items {
		ic, ok := c.Get(ctx, i.Path)
		require.Zero(t, ic)
		require.False(t, ok)
	}
	require.Equal(t, 4, c.Len())
	require.Equal(t, 16, c.Size())

	c.Set(ctx, PreRenderedItem{
		Path: "/test5",
		Body: []byte("test"),
	})
	require.Equal(t, 1, c.Len())
	require.Equal(t, 4, c.Size())
	require.False(t, evictCalled)
}

func TestPreRenderLRUCacheEvict(t *testing.T) {
	ctx := context.TODO()

	evictCount := 0
	evictSize := 0
	onEvict := func(path string, i PreRenderedItem) {
		evictCount++
		evictSize += i.Size()
	}

	c := NewPreRenderLRUCache(8, time.Second, onEvict).(*preRenderLRUCache)

	items := []PreRenderedItem{
		{
			Path: "/test1",
			Body: []byte("test"),
		},
		{
			Path: "/test2",
			Body: []byte("test"),
		},
	}

	for _, i := range items {
		c.Set(ctx, i)
	}
	require.Equal(t, 2, c.Len())
	require.Equal(t, 8, c.Size())
	require.Equal(t, 0, evictCount)
	require.Equal(t, 0, evictSize)

	c.Set(ctx, PreRenderedItem{
		Path: "/test3",
		Body: []byte("test"),
	})
	require.Equal(t, 2, c.Len())
	require.Equal(t, 8, c.Size())
	require.Equal(t, 1, evictCount)
	require.Equal(t, 4, evictSize)

	c.Set(ctx, PreRenderedItem{
		Path: "/test4",
		Body: []byte("testbig"),
	})
	require.Equal(t, 1, c.Len())
	require.Equal(t, 7, c.Size())
	require.Equal(t, 3, evictCount)
	require.Equal(t, 12, evictSize)
}
