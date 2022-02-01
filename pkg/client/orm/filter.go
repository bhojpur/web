package orm

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
)

// FilterChain is used to build a Filter
// don't forget to call next(...) inside your Filter
type FilterChain func(next Filter) Filter

// Filter's behavior is a little big strange.
// it's only be called when users call methods of Ormer
// return value is an array. it's a little bit hard to understand,
// for example, the Ormer's Read method only return error
// so the filter processing this method should return an array whose first element is error
// and, Ormer's ReadOrCreateWithCtx return three values, so the Filter's result should contains three values
type Filter func(ctx context.Context, inv *Invocation) []interface{}

var globalFilterChains = make([]FilterChain, 0, 4)

// AddGlobalFilterChain adds a new FilterChain
// All orm instances built after this invocation will use this filterChain,
// but instances built before this invocation will not be affected
func AddGlobalFilterChain(filterChain ...FilterChain) {
	globalFilterChains = append(globalFilterChains, filterChain...)
}
