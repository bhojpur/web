package metric

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
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bhojpur/web/pkg/adapter/context"
)

func TestPrometheusMiddleWare(t *testing.T) {
	middleware := PrometheusMiddleWare(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	writer := &context.Response{}
	request := &http.Request{
		URL: &url.URL{
			Host:    "localhost",
			RawPath: "/a/b/c",
		},
		Method: "POST",
	}
	vec := prometheus.NewSummaryVec(prometheus.SummaryOpts{}, []string{"pattern", "method", "status", "duration"})

	report(time.Second, writer, request, vec)
	middleware.ServeHTTP(writer, request)
}
