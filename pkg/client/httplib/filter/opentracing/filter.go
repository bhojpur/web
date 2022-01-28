package opentracing

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
