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
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	webapp "github.com/bhojpur/web/pkg/adapter"
	ctxsvr "github.com/bhojpur/web/pkg/context"
	websvr "github.com/bhojpur/web/pkg/engine"
)

// FilterChainBuilder is an extension point,
// when we want to support some configuration,
// please use this structure
type FilterChainBuilder struct {
}

// FilterChain returns a FilterFunc. The filter will records some metrics
func (builder *FilterChainBuilder) FilterChain(next websvr.FilterFunc) websvr.FilterFunc {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "bhojpur",
		Subsystem: "http_request",
		ConstLabels: map[string]string{
			"server":  websvr.BConfig.ServerName,
			"env":     websvr.BConfig.RunMode,
			"appname": websvr.BConfig.AppName,
		},
		Help: "The statistics info for HTTP request",
	}, []string{"pattern", "method", "status", "duration"})

	prometheus.MustRegister(summaryVec)

	registerBuildInfo()

	return func(ctx *ctxsvr.Context) {
		startTime := time.Now()
		next(ctx)
		endTime := time.Now()
		go report(endTime.Sub(startTime), ctx, summaryVec)
	}
}

func registerBuildInfo() {
	buildInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name:      "bhojpur",
		Subsystem: "build_info",
		Help:      "The building information",
		ConstLabels: map[string]string{
			"appname":        websvr.BConfig.AppName,
			"build_version":  webapp.BuildVersion,
			"build_revision": webapp.BuildGitRevision,
			"build_status":   webapp.BuildStatus,
			"build_tag":      webapp.BuildTag,
			"build_time":     strings.Replace(webapp.BuildTime, "--", " ", 1),
			"go_version":     webapp.GoVersion,
			"git_branch":     webapp.GitBranch,
			"start_time":     time.Now().Format("2018-03-26 15:04:05"),
		},
	}, []string{})

	prometheus.MustRegister(buildInfo)
	buildInfo.WithLabelValues().Set(1)
}

func report(dur time.Duration, ctx *ctxsvr.Context, vec *prometheus.SummaryVec) {
	status := ctx.Output.Status
	ptn := ctx.Input.GetData("RouterPattern").(string)
	ms := dur / time.Millisecond
	vec.WithLabelValues(ptn, ctx.Input.Method(), strconv.Itoa(status), strconv.Itoa(int(ms))).Observe(float64(ms))
}
