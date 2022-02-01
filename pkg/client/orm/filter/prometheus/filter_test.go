package prometheus

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

	"github.com/stretchr/testify/assert"

	"github.com/bhojpur/web/pkg/client/orm"
)

func TestFilterChainBuilder_FilterChain1(t *testing.T) {
	next := func(ctx context.Context, inv *orm.Invocation) []interface{} {
		inv.Method = "coming"
		return []interface{}{}
	}
	builder := &FilterChainBuilder{}
	filter := builder.FilterChain(next)

	assert.NotNil(t, builder.summaryVec)
	assert.NotNil(t, filter)

	inv := &orm.Invocation{}
	filter(context.Background(), inv)
	assert.Equal(t, "coming", inv.Method)

	inv = &orm.Invocation{
		Method:      "Hello",
		TxStartTime: time.Now(),
	}
	builder.reportTxn(context.Background(), inv)

	inv = &orm.Invocation{
		Method: "Begin",
	}

	ctx := context.Background()
	// it will be ignored
	builder.report(ctx, inv, time.Second)

	inv.Method = "Commit"
	builder.report(ctx, inv, time.Second)

	inv.Method = "Update"
	builder.report(ctx, inv, time.Second)

}
