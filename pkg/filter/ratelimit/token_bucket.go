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
	"sync"
	"time"
)

type tokenBucket struct {
	sync.RWMutex
	remaining   uint
	capacity    uint
	lastCheckAt time.Time
	rate        time.Duration
}

// newTokenBucket return an bucket that implements token bucket
func newTokenBucket(opts ...bucketOption) bucket {
	b := &tokenBucket{lastCheckAt: time.Now()}
	for _, o := range opts {
		o(b)
	}
	return b
}

func withCapacity(capacity uint) bucketOption {
	return func(b bucket) {
		bucket := b.(*tokenBucket)
		bucket.capacity = capacity
		bucket.remaining = capacity
	}
}

func withRate(rate time.Duration) bucketOption {
	return func(b bucket) {
		bucket := b.(*tokenBucket)
		bucket.rate = rate
	}
}

func (b *tokenBucket) getRemaining() uint {
	b.RLock()
	defer b.RUnlock()
	return b.remaining
}

func (b *tokenBucket) getRate() time.Duration {
	b.RLock()
	defer b.RUnlock()
	return b.rate
}

func (b *tokenBucket) getCapacity() uint {
	b.RLock()
	defer b.RUnlock()
	return b.capacity
}

func (b *tokenBucket) take(amount uint) bool {
	if b.rate <= 0 {
		return true
	}
	b.Lock()
	defer b.Unlock()
	now := time.Now()
	times := uint(now.Sub(b.lastCheckAt) / b.rate)
	b.lastCheckAt = b.lastCheckAt.Add(time.Duration(times) * b.rate)
	b.remaining += times
	if b.remaining < amount {
		return false
	}
	b.remaining -= amount
	if b.remaining > b.capacity {
		b.remaining = b.capacity
	}
	return true
}
