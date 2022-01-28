package adapter

import (
	"github.com/bhojpur/web/pkg/adapter/context"
	bcontext "github.com/bhojpur/web/pkg/context"
	web "github.com/bhojpur/web/pkg/engine"
)

// FilterFunc defines a filter function which is invoked before the controller handler is executed.
type FilterFunc func(*context.Context)

// FilterRouter defines a filter operation which is invoked before the controller handler is executed.
// It can match the URL against a pattern, and execute a filter function
// when a request with a matching URL arrives.
type FilterRouter web.FilterRouter

// ValidRouter checks if the current request is matched by this filter.
// If the request is matched, the values of the URL parameters defined
// by the filter pattern are also returned.
func (f *FilterRouter) ValidRouter(url string, ctx *context.Context) bool {
	return (*web.FilterRouter)(f).ValidRouter(url, (*bcontext.Context)(ctx))
}
