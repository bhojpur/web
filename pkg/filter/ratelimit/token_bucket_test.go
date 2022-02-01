package ratelimit

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
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	b := newTokenBucket(withRate(1 * time.Second)).(*tokenBucket)
	assert.Equal(t, b.getRate(), 1*time.Second)
}

func TestGetRemainingAndCapacity(t *testing.T) {
	b := newTokenBucket(withCapacity(10))
	assert.Equal(t, b.getRemaining(), uint(10))
	assert.Equal(t, b.getCapacity(), uint(10))
}

func TestTake(t *testing.T) {
	b := newTokenBucket(withCapacity(10), withRate(10*time.Millisecond)).(*tokenBucket)
	for i := 0; i < 10; i++ {
		assert.True(t, b.take(1))
	}
	assert.False(t, b.take(1))
	assert.Equal(t, b.getRemaining(), uint(0))
	b = newTokenBucket(withCapacity(1), withRate(1*time.Millisecond)).(*tokenBucket)
	assert.True(t, b.take(1))
	time.Sleep(2 * time.Millisecond)
	assert.True(t, b.take(1))
}
