package log

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
	"io"
	"net/http"
	"net/http/httputil"

	logs "github.com/bhojpur/logger/pkg/engine"
	"github.com/bhojpur/web/pkg/client/httplib"
)

// FilterChainBuilder can build a log filter
type FilterChainBuilder struct {
	printableContentTypes []string                              // only print the body of included mime types of request and response
	log                   func(f interface{}, v ...interface{}) // custom log function
}

// BuilderOption option constructor
type BuilderOption func(*FilterChainBuilder)

type logInfo struct {
	req  []byte
	resp []byte
	err  error
}

var defaultprintableContentTypes = []string{
	"text/plain", "text/xml", "text/html", "text/csv",
	"text/calendar", "text/javascript", "text/javascript",
	"text/css",
}

// NewFilterChainBuilder initialize a filterChainBuilder, pass options to customize
func NewFilterChainBuilder(opts ...BuilderOption) *FilterChainBuilder {
	res := &FilterChainBuilder{
		printableContentTypes: defaultprintableContentTypes,
		log:                   logs.Debug,
	}
	for _, o := range opts {
		o(res)
	}

	return res
}

// WithLog return option constructor modify log function
func WithLog(f func(f interface{}, v ...interface{})) BuilderOption {
	return func(h *FilterChainBuilder) {
		h.log = f
	}
}

// WithprintableContentTypes return option constructor modify printableContentTypes
func WithprintableContentTypes(types []string) BuilderOption {
	return func(h *FilterChainBuilder) {
		h.printableContentTypes = types
	}
}

// FilterChain can print the request after FilterChain processing and response before processsing
func (builder *FilterChainBuilder) FilterChain(next httplib.Filter) httplib.Filter {
	return func(ctx context.Context, req *httplib.BhojpurHTTPRequest) (*http.Response, error) {
		info := &logInfo{}
		defer info.print(builder.log)
		resp, err := next(ctx, req)
		info.err = err
		contentType := req.GetRequest().Header.Get("Content-Type")
		shouldPrintBody := builder.shouldPrintBody(contentType, req.GetRequest().Body)
		dump, err := httputil.DumpRequest(req.GetRequest(), shouldPrintBody)
		info.req = dump
		if err != nil {
			logs.Error(err)
		}
		if resp != nil {
			contentType = resp.Header.Get("Content-Type")
			shouldPrintBody = builder.shouldPrintBody(contentType, resp.Body)
			dump, err = httputil.DumpResponse(resp, shouldPrintBody)
			info.resp = dump
			if err != nil {
				logs.Error(err)
			}
		}
		return resp, err
	}
}

func (builder *FilterChainBuilder) shouldPrintBody(contentType string, body io.ReadCloser) bool {
	if contains(builder.printableContentTypes, contentType) {
		return true
	}
	if body != nil {
		logs.Warn("printableContentTypes do not contain %s, if you want to print request and response body please add it.", contentType)
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (info *logInfo) print(log func(f interface{}, v ...interface{})) {
	log("Request: ====================")
	log("%q", info.req)
	log("Response: ===================")
	log("%q", info.resp)
	if info.err != nil {
		log("Error: ======================")
		log("%q", info.err)
	}
}
