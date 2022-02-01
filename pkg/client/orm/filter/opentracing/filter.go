package opentracing

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
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/bhojpur/web/pkg/client/orm"
)

// FilterChainBuilder provides an extension point
// this Filter's behavior looks a little bit strange
// for example:
// if we want to trace QuerySetter
// actually we trace invoking "QueryTable" and "QueryTableWithCtx"
// the method Begin*, Commit and Rollback are ignored.
// When use using those methods, it means that they want to manager their transaction manually, so we won't handle them.
type FilterChainBuilder struct {
	// CustomSpanFunc users are able to custom their span
	CustomSpanFunc func(span opentracing.Span, ctx context.Context, inv *orm.Invocation)
}

func (builder *FilterChainBuilder) FilterChain(next orm.Filter) orm.Filter {
	return func(ctx context.Context, inv *orm.Invocation) []interface{} {
		operationName := builder.operationName(ctx, inv)
		if strings.HasPrefix(inv.Method, "Begin") || inv.Method == "Commit" || inv.Method == "Rollback" {
			return next(ctx, inv)
		}

		span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()
		res := next(spanCtx, inv)
		builder.buildSpan(span, spanCtx, inv)
		return res
	}
}

func (builder *FilterChainBuilder) buildSpan(span opentracing.Span, ctx context.Context, inv *orm.Invocation) {
	span.SetTag("orm.method", inv.Method)
	span.SetTag("orm.table", inv.GetTableName())
	span.SetTag("orm.insideTx", inv.InsideTx)
	span.SetTag("orm.txName", ctx.Value(orm.TxNameKey))
	span.SetTag("span.kind", "client")
	span.SetTag("component", "bhojpur")

	if builder.CustomSpanFunc != nil {
		builder.CustomSpanFunc(span, ctx, inv)
	}
}

func (builder *FilterChainBuilder) operationName(ctx context.Context, inv *orm.Invocation) string {
	if n, ok := ctx.Value(orm.TxNameKey).(string); ok {
		return inv.Method + "#tx(" + n + ")"
	}
	return inv.Method + "#" + inv.GetTableName()
}
