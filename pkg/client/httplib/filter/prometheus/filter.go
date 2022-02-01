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
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bhojpur/web/pkg/client/httplib"
)

type FilterChainBuilder struct {
	summaryVec prometheus.ObserverVec
	AppName    string
	ServerName string
	RunMode    string
}

func (builder *FilterChainBuilder) FilterChain(next httplib.Filter) httplib.Filter {

	builder.summaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "bhojpur",
		Subsystem: "remote_http_request",
		ConstLabels: map[string]string{
			"server":  builder.ServerName,
			"env":     builder.RunMode,
			"appname": builder.AppName,
		},
		Help: "The statics info for remote http requests",
	}, []string{"proto", "scheme", "method", "host", "path", "status", "duration", "isError"})

	return func(ctx context.Context, req *httplib.BhojpurHTTPRequest) (*http.Response, error) {
		startTime := time.Now()
		resp, err := next(ctx, req)
		endTime := time.Now()
		go builder.report(startTime, endTime, ctx, req, resp, err)
		return resp, err
	}
}

func (builder *FilterChainBuilder) report(startTime time.Time, endTime time.Time,
	ctx context.Context, req *httplib.BhojpurHTTPRequest, resp *http.Response, err error) {

	proto := req.GetRequest().Proto

	scheme := req.GetRequest().URL.Scheme
	method := req.GetRequest().Method

	host := req.GetRequest().URL.Host
	path := req.GetRequest().URL.Path

	status := -1
	if resp != nil {
		status = resp.StatusCode
	}

	dur := int(endTime.Sub(startTime) / time.Millisecond)

	builder.summaryVec.WithLabelValues(proto, scheme, method, host, path,
		strconv.Itoa(status), strconv.Itoa(dur), strconv.FormatBool(err == nil))
}
