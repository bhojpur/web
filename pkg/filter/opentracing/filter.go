package opentracing

import (
	"context"

	ctxutil "github.com/bhojpur/web/pkg/context"
	web "github.com/bhojpur/web/pkg/engine"
	logKit "github.com/go-kit/kit/log"
	opentracingKit "github.com/go-kit/kit/tracing/opentracing"
	"github.com/opentracing/opentracing-go"
)

// FilterChainBuilder provides an extension point that we can support more configurations if necessary
type FilterChainBuilder struct {
	// CustomSpanFunc makes users to custom the span.
	CustomSpanFunc func(span opentracing.Span, ctx *ctxutil.Context)
}

func (builder *FilterChainBuilder) FilterChain(next web.FilterFunc) web.FilterFunc {
	return func(ctx *ctxutil.Context) {
		var (
			spanCtx context.Context
			span    opentracing.Span
		)
		operationName := builder.operationName(ctx)

		if preSpan := opentracing.SpanFromContext(ctx.Request.Context()); preSpan == nil {
			inject := opentracingKit.HTTPToContext(opentracing.GlobalTracer(), operationName, logKit.NewNopLogger())
			spanCtx = inject(ctx.Request.Context(), ctx.Request)
			span = opentracing.SpanFromContext(spanCtx)
		} else {
			span, spanCtx = opentracing.StartSpanFromContext(ctx.Request.Context(), operationName)
		}

		defer span.Finish()

		newReq := ctx.Request.Clone(spanCtx)
		ctx.Reset(ctx.ResponseWriter.ResponseWriter, newReq)

		next(ctx)
		// if you think we need to do more things, feel free to create an issue to tell us
		span.SetTag("http.status_code", ctx.ResponseWriter.Status)
		span.SetTag("http.method", ctx.Input.Method())
		span.SetTag("peer.hostname", ctx.Request.Host)
		span.SetTag("http.url", ctx.Request.URL.String())
		span.SetTag("http.scheme", ctx.Request.URL.Scheme)
		span.SetTag("span.kind", "server")
		span.SetTag("component", "bhojpur")
		if ctx.Output.IsServerError() || ctx.Output.IsClientError() {
			span.SetTag("error", true)
		}
		span.SetTag("peer.address", ctx.Request.RemoteAddr)
		span.SetTag("http.proto", ctx.Request.Proto)

		span.SetTag("bhojpur.route", ctx.Input.GetData("RouterPattern"))

		if builder.CustomSpanFunc != nil {
			builder.CustomSpanFunc(span, ctx)
		}
	}
}

func (builder *FilterChainBuilder) operationName(ctx *ctxutil.Context) string {
	operationName := ctx.Input.URL()
	// it means that there is not any span, so we create a span as the root span.
	// TODO, if we support multiple servers, this need to be changed
	route, found := web.BhojpurApp.Handlers.FindRouter(ctx)
	if found {
		operationName = ctx.Input.Method() + "#" + route.GetPattern()
	}
	return operationName
}
