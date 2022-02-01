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

	"github.com/bhojpur/web/pkg/context"
	web "github.com/bhojpur/web/pkg/engine"
)

// limiterOption is constructor option
type limiterOption func(l *limiter)

type limiter struct {
	sync.RWMutex
	capacity      uint
	rate          time.Duration
	buckets       map[string]bucket
	bucketFactory func(opts ...bucketOption) bucket
	sessionKey    func(ctx *context.Context) string
	resp          RejectionResponse
}

// RejectionResponse stores response information
// for the request rejected by limiter
type RejectionResponse struct {
	code int
	body string
}

const perRequestConsumedAmount = 1

var defaultRejectionResponse = RejectionResponse{
	code: 429,
	body: "too many requests",
}

// NewLimiter return FilterFunc, the limiter enables rate limit
// according to the configuration.
func NewLimiter(opts ...limiterOption) web.FilterFunc {
	l := &limiter{
		buckets:       make(map[string]bucket),
		sessionKey:    defaultSessionKey,
		rate:          time.Millisecond * 10,
		capacity:      100,
		bucketFactory: newTokenBucket,
		resp:          defaultRejectionResponse,
	}
	for _, o := range opts {
		o(l)
	}

	return func(ctx *context.Context) {
		if !l.take(perRequestConsumedAmount, ctx) {
			ctx.ResponseWriter.WriteHeader(l.resp.code)
			ctx.WriteString(l.resp.body)
		}
	}
}

// WithSessionKey return limiterOption. WithSessionKey config func
// which defines the request characteristic against the limit is applied
func WithSessionKey(f func(ctx *context.Context) string) limiterOption {
	return func(l *limiter) {
		l.sessionKey = f
	}
}

// WithRate return limiterOption. WithRate config how long it takes to
// generate a token.
func WithRate(r time.Duration) limiterOption {
	return func(l *limiter) {
		l.rate = r
	}
}

// WithCapacity return limiterOption. WithCapacity config the capacity size.
// The bucket with a capacity of n has n tokens after initialization. The capacity
// defines how many requests a client can make in excess of the rate.
func WithCapacity(c uint) limiterOption {
	return func(l *limiter) {
		l.capacity = c
	}
}

// WithBucketFactory return limiterOption. WithBucketFactory customize the
// implementation of Bucket.
func WithBucketFactory(f func(opts ...bucketOption) bucket) limiterOption {
	return func(l *limiter) {
		l.bucketFactory = f
	}
}

// WithRejectionResponse return limiterOption. WithRejectionResponse
// customize the response for the request rejected by the limiter.
func WithRejectionResponse(resp RejectionResponse) limiterOption {
	return func(l *limiter) {
		l.resp = resp
	}
}

func (l *limiter) take(amount uint, ctx *context.Context) bool {
	bucket := l.getBucket(ctx)
	if bucket == nil {
		return true
	}
	return bucket.take(amount)
}

func (l *limiter) getBucket(ctx *context.Context) bucket {
	key := l.sessionKey(ctx)
	l.RLock()
	b, ok := l.buckets[key]
	l.RUnlock()
	if !ok {
		b = l.createBucket(key)
	}

	return b
}

func (l *limiter) createBucket(key string) bucket {
	l.Lock()
	defer l.Unlock()
	// double check avoid overwriting
	b, ok := l.buckets[key]
	if ok {
		return b
	}
	b = l.bucketFactory(withCapacity(l.capacity), withRate(l.rate))
	l.buckets[key] = b
	return b
}

func defaultSessionKey(ctx *context.Context) string {
	return "BHOJPUR_ALL"
}

func RemoteIPSessionKey(ctx *context.Context) string {
	r := ctx.Request
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
