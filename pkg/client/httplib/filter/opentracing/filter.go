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
	"net/http"

	"github.com/bhojpur/web/pkg/client/httplib"
	logKit "github.com/go-kit/kit/log"
	opentracingKit "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
)

type FilterChainBuilder struct {
	// CustomSpanFunc users are able to custom their span
	CustomSpanFunc func(span opentracing.Span, ctx context.Context,
		req *httplib.BhojpurHTTPRequest, resp *http.Response, err error)
}

func (builder *FilterChainBuilder) FilterChain(next httplib.Filter) httplib.Filter {

	return func(ctx context.Context, req *httplib.BhojpurHTTPRequest) (*http.Response, error) {

		method := req.GetRequest().Method

		operationName := method + "#" + req.GetRequest().URL.String()
		span, spanCtx := opentracing.StartSpanFromContext(ctx, operationName)
		defer span.Finish()

		inject := opentracingKit.ContextToHTTP(opentracing.GlobalTracer(), logKit.NewNopLogger())
		inject(spanCtx, req.GetRequest())
		resp, err := next(spanCtx, req)

		if resp != nil {
			span.SetTag("http.status_code", resp.StatusCode)
		}
		span.SetTag("http.method", method)
		span.SetTag("peer.hostname", req.GetRequest().URL.Host)
		span.SetTag("http.url", req.GetRequest().URL.String())
		span.SetTag("http.scheme", req.GetRequest().URL.Scheme)
		span.SetTag("span.kind", "client")
		span.SetTag("component", "bhojpur")
		if err != nil {
			span.SetTag("error", true)
			span.SetTag("message", err.Error())
		} else if resp != nil && !(resp.StatusCode < 300 && resp.StatusCode >= 200) {
			span.SetTag("error", true)
		}

		span.SetTag("peer.address", req.GetRequest().RemoteAddr)
		span.SetTag("http.proto", req.GetRequest().Proto)

		if builder.CustomSpanFunc != nil {
			builder.CustomSpanFunc(span, ctx, req, resp, err)
		}
		return resp, err
	}
}
