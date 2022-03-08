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
	"time"

	"github.com/stretchr/testify/require"
)

func TestExpire(t *testing.T) {
	testCache(t, &Expire{
		ItemTTL: time.Minute,
	})
}

func TestExpireExpire(t *testing.T) {
	ctx := context.TODO()

	c := Expire{
		ItemTTL: time.Minute,
	}

	c.Set(ctx, "/hello", String("hello"))
	require.Equal(t, 1, c.Len())
	require.Equal(t, 5, c.Size())

	c.Set(ctx, "/world", String("world"))
	require.Equal(t, 2, c.Len())
	require.Equal(t, 10, c.Size())

	c.Set(ctx, "/goodbye", String("goodbye"))
	require.Equal(t, 3, c.Len())
	require.Equal(t, 17, c.Size())

	c.items["/hello"].expiresAt = time.Now().Add(-time.Second)
	c.items["/world"].expiresAt = time.Now().Add(-time.Second)

	c.Set(ctx, "/goodmorning", String("goodmorning"))
	require.Equal(t, 2, c.Len())
	require.Equal(t, 18, c.Size())
	require.Len(t, c.queue, 2)
}

func TestExpireSetSameKey(t *testing.T) {
	ctx := context.TODO()

	c := Expire{
		ItemTTL: time.Minute,
	}

	c.Set(ctx, "/foo", String("foo"))
	require.Len(t, c.queue, 1)
	require.Equal(t, 1, c.Len())
	require.Equal(t, 3, c.Size())

	c.Set(ctx, "/bar", String("bar"))
	require.Len(t, c.queue, 2)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 6, c.Size())

	c.Set(ctx, "/bar", String("barre"))
	require.Len(t, c.queue, 3)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 8, c.Size())

	c.Set(ctx, "/foo", String("fooo"))
	require.Len(t, c.queue, 2)
	require.Equal(t, 2, c.Len())
	require.Equal(t, 9, c.Size())
}
