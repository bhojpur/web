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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	logs "github.com/bhojpur/logger/pkg/engine"
	bhojpur "github.com/bhojpur/web/pkg/engine"
	web "github.com/bhojpur/web/pkg/engine"
)

func PrometheusMiddleWare(next http.Handler) http.Handler {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "bhojpur",
		Subsystem: "http_request",
		ConstLabels: map[string]string{
			"server":  web.BConfig.ServerName,
			"env":     web.BConfig.RunMode,
			"appname": web.BConfig.AppName,
		},
		Help: "The statics info for http request",
	}, []string{"pattern", "method", "status", "duration"})

	prometheus.MustRegister(summaryVec)

	registerBuildInfo()

	return http.HandlerFunc(func(writer http.ResponseWriter, q *http.Request) {
		start := time.Now()
		next.ServeHTTP(writer, q)
		end := time.Now()
		go report(end.Sub(start), writer, q, summaryVec)
	})
}

func registerBuildInfo() {
	buildInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "bhojpur",
		Subsystem: "build_info",
		Help:      "The building information",
		ConstLabels: map[string]string{
			"appname":        web.BConfig.AppName,
			"build_version":  bhojpur.BuildVersion,
			"build_revision": bhojpur.BuildGitRevision,
			"build_status":   bhojpur.BuildStatus,
			"build_tag":      bhojpur.BuildTag,
			"build_time":     strings.Replace(bhojpur.BuildTime, "--", " ", 1),
			"go_version":     bhojpur.GoVersion,
			"git_branch":     bhojpur.GitBranch,
			"start_time":     time.Now().Format("2006-01-02 15:04:05"),
		},
	}, []string{})

	prometheus.MustRegister(buildInfo)
	buildInfo.WithLabelValues().Set(1)
}

func report(dur time.Duration, writer http.ResponseWriter, q *http.Request, vec *prometheus.SummaryVec) {
	ctrl := web.BhojpurApp.Handlers
	ctx := ctrl.GetContext()
	ctx.Reset(writer, q)
	defer ctrl.GiveBackContext(ctx)

	// We cannot read the status code from q.Response.StatusCode
	// since the http server does not set q.Response. So q.Response is nil
	// Thus, we use reflection to read the status from writer whose concrete type is http.response
	responseVal := reflect.ValueOf(writer).Elem()
	field := responseVal.FieldByName("status")
	status := -1
	if field.IsValid() && field.Kind() == reflect.Int {
		status = int(field.Int())
	}
	ptn := "UNKNOWN"
	if rt, found := ctrl.FindRouter(ctx); found {
		ptn = rt.GetPattern()
	} else {
		logs.Warn("we can not find the router info for this request, so request will be recorded as UNKNOWN: " + q.URL.String())
	}
	ms := dur / time.Millisecond
	vec.WithLabelValues(ptn, q.Method, strconv.Itoa(status), strconv.Itoa(int(ms))).Observe(float64(ms))
}
