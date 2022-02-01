package cors

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

// Package cors provides handlers to enable CORS support.
// Usage
//	import (
// 		bhojpur "github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/plugins/cors"
// )
//
//	func main() {
//		// CORS for https://foo.* origins, allowing:
//		// - PUT and PATCH methods
//		// - Origin header
//		// - Credentials share
//		bhojpur.InsertFilter("*", bhojpur.BeforeRouter, cors.Allow(&cors.Options{
//			AllowOrigins:     []string{"https://*.foo.com"},
//			AllowMethods:     []string{"PUT", "PATCH"},
//			AllowHeaders:     []string{"Origin"},
//			ExposeHeaders:    []string{"Content-Length"},
//			AllowCredentials: true,
//		}))
//		bhojpur.Run()
//	}

import (
	bhojpur "github.com/bhojpur/web/pkg/adapter"
	beecontext "github.com/bhojpur/web/pkg/context"
	"github.com/bhojpur/web/pkg/filter/cors"

	"github.com/bhojpur/web/pkg/adapter/context"
)

// Options represents Access Control options.
type Options cors.Options

// Header converts options into CORS headers.
func (o *Options) Header(origin string) (headers map[string]string) {
	return (*cors.Options)(o).Header(origin)
}

// PreflightHeader converts options into CORS headers for a preflight response.
func (o *Options) PreflightHeader(origin, rMethod, rHeaders string) (headers map[string]string) {
	return (*cors.Options)(o).PreflightHeader(origin, rMethod, rHeaders)
}

// IsOriginAllowed looks up if the origin matches one of the patterns
// generated from Options.AllowOrigins patterns.
func (o *Options) IsOriginAllowed(origin string) bool {
	return (*cors.Options)(o).IsOriginAllowed(origin)
}

// Allow enables CORS for requests those match the provided options.
func Allow(opts *Options) bhojpur.FilterFunc {
	f := cors.Allow((*cors.Options)(opts))
	return func(c *context.Context) {
		f((*beecontext.Context)(c))
	}
}
