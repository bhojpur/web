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

// It provides handlers to enable CORS support.
// Usage
//	import (
// 		websvr "github.com/bhojpur/web/pkg/engine"
//		"github.com/bhojpur/web/pkg/filter/plugins/cors"
// )
//
//	func main() {
//		// CORS for https://foo.* origins, allowing:
//		// - PUT and PATCH methods
//		// - Origin header
//		// - Credentials share
//		websvr.InsertFilter("*", websvr.BeforeRouter, cors.Allow(&cors.Options{
//			AllowOrigins:     []string{"https://*.bhojpur.net"},
//			AllowMethods:     []string{"PUT", "PATCH"},
//			AllowHeaders:     []string{"Origin"},
//			ExposeHeaders:    []string{"Content-Length"},
//			AllowCredentials: true,
//		}))
//		websvr.Run()
//	}

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

const (
	headerAllowOrigin      = "Access-Control-Allow-Origin"
	headerAllowCredentials = "Access-Control-Allow-Credentials"
	headerAllowHeaders     = "Access-Control-Allow-Headers"
	headerAllowMethods     = "Access-Control-Allow-Methods"
	headerExposeHeaders    = "Access-Control-Expose-Headers"
	headerMaxAge           = "Access-Control-Max-Age"

	headerOrigin         = "Origin"
	headerRequestMethod  = "Access-Control-Request-Method"
	headerRequestHeaders = "Access-Control-Request-Headers"
)

var (
	defaultAllowHeaders = []string{"Origin", "Accept", "Content-Type", "Authorization"}
	// Regex patterns are generated from AllowOrigins. These are used and generated internally.
	allowOriginPatterns = []string{}
)

// Options represents Access Control options.
type Options struct {
	// If set, all origins are allowed.
	AllowAllOrigins bool
	// A list of allowed origins. Wild cards and FQDNs are supported.
	AllowOrigins []string
	// If set, allows to share auth credentials such as cookies.
	AllowCredentials bool
	// A list of allowed HTTP methods.
	AllowMethods []string
	// A list of allowed HTTP headers.
	AllowHeaders []string
	// A list of exposed HTTP headers.
	ExposeHeaders []string
	// Max age of the CORS headers.
	MaxAge time.Duration
}

// Header converts options into CORS headers.
func (o *Options) Header(origin string) (headers map[string]string) {
	headers = make(map[string]string)
	// if origin is not allowed, don't extend the headers
	// with CORS headers.
	if !o.AllowAllOrigins && !o.IsOriginAllowed(origin) {
		return
	}

	// add allow origin
	if o.AllowAllOrigins {
		headers[headerAllowOrigin] = "*"
	} else {
		headers[headerAllowOrigin] = origin
	}

	// add allow credentials
	headers[headerAllowCredentials] = strconv.FormatBool(o.AllowCredentials)

	// add allow methods
	if len(o.AllowMethods) > 0 {
		headers[headerAllowMethods] = strings.Join(o.AllowMethods, ",")
	}

	// add allow headers
	if len(o.AllowHeaders) > 0 {
		headers[headerAllowHeaders] = strings.Join(o.AllowHeaders, ",")
	}

	// add exposed header
	if len(o.ExposeHeaders) > 0 {
		headers[headerExposeHeaders] = strings.Join(o.ExposeHeaders, ",")
	}
	// add a max age header
	if o.MaxAge > time.Duration(0) {
		headers[headerMaxAge] = strconv.FormatInt(int64(o.MaxAge/time.Second), 10)
	}
	return
}

// PreflightHeader converts options into CORS headers for a preflight response.
func (o *Options) PreflightHeader(origin, rMethod, rHeaders string) (headers map[string]string) {
	headers = make(map[string]string)
	if !o.AllowAllOrigins && !o.IsOriginAllowed(origin) {
		return
	}
	// verify if requested method is allowed
	for _, method := range o.AllowMethods {
		if method == rMethod {
			headers[headerAllowMethods] = strings.Join(o.AllowMethods, ",")
			break
		}
	}

	// verify if requested headers are allowed
	var allowed []string
	for _, rHeader := range strings.Split(rHeaders, ",") {
		rHeader = strings.TrimSpace(rHeader)
	lookupLoop:
		for _, allowedHeader := range o.AllowHeaders {
			if strings.ToLower(rHeader) == strings.ToLower(allowedHeader) {
				allowed = append(allowed, rHeader)
				break lookupLoop
			}
		}
	}

	headers[headerAllowCredentials] = strconv.FormatBool(o.AllowCredentials)
	// add allow origin
	if o.AllowAllOrigins {
		headers[headerAllowOrigin] = "*"
	} else {
		headers[headerAllowOrigin] = origin
	}

	// add allowed headers
	if len(allowed) > 0 {
		headers[headerAllowHeaders] = strings.Join(allowed, ",")
	}

	// add exposed headers
	if len(o.ExposeHeaders) > 0 {
		headers[headerExposeHeaders] = strings.Join(o.ExposeHeaders, ",")
	}
	// add a max age header
	if o.MaxAge > time.Duration(0) {
		headers[headerMaxAge] = strconv.FormatInt(int64(o.MaxAge/time.Second), 10)
	}
	return
}

// IsOriginAllowed looks up if the origin matches one of the patterns
// generated from Options.AllowOrigins patterns.
func (o *Options) IsOriginAllowed(origin string) (allowed bool) {
	for _, pattern := range allowOriginPatterns {
		allowed, _ = regexp.MatchString(pattern, origin)
		if allowed {
			return
		}
	}
	return
}

// Allow enables CORS for requests those match the provided options.
func Allow(opts *Options) websvr.FilterFunc {
	// Allow default headers if nothing is specified.
	if len(opts.AllowHeaders) == 0 {
		opts.AllowHeaders = defaultAllowHeaders
	}

	for _, origin := range opts.AllowOrigins {
		pattern := regexp.QuoteMeta(origin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		allowOriginPatterns = append(allowOriginPatterns, "^"+pattern+"$")
	}

	return func(ctx *context.Context) {
		var (
			origin           = ctx.Input.Header(headerOrigin)
			requestedMethod  = ctx.Input.Header(headerRequestMethod)
			requestedHeaders = ctx.Input.Header(headerRequestHeaders)
			// additional headers to be added
			// to the response.
			headers map[string]string
		)

		if ctx.Input.Method() == "OPTIONS" &&
			(requestedMethod != "" || requestedHeaders != "") {
			headers = opts.PreflightHeader(origin, requestedMethod, requestedHeaders)
			for key, value := range headers {
				ctx.Output.Header(key, value)
			}
			ctx.ResponseWriter.WriteHeader(http.StatusOK)
			return
		}
		headers = opts.Header(origin)

		for key, value := range headers {
			ctx.Output.Header(key, value)
		}
	}
}
