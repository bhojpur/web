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
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bhojpur/web/pkg/client/orm"
)

// FilterChainBuilder is an extension point,
// when we want to support some configuration,
// please use this structure
// this Filter's behavior looks a little bit strange
// for example:
// if we want to records the metrics of QuerySetter
// actually we only records metrics of invoking "QueryTable" and "QueryTableWithCtx"
type FilterChainBuilder struct {
	summaryVec prometheus.ObserverVec
	AppName    string
	ServerName string
	RunMode    string
}

func (builder *FilterChainBuilder) FilterChain(next orm.Filter) orm.Filter {

	builder.summaryVec = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "bhojpur",
		Subsystem: "orm_operation",
		ConstLabels: map[string]string{
			"server":  builder.ServerName,
			"env":     builder.RunMode,
			"appname": builder.AppName,
		},
		Help: "The statics info for orm operation",
	}, []string{"method", "name", "duration", "insideTx", "txName"})

	return func(ctx context.Context, inv *orm.Invocation) []interface{} {
		startTime := time.Now()
		res := next(ctx, inv)
		endTime := time.Now()
		dur := (endTime.Sub(startTime)) / time.Millisecond

		// if the TPS is too large, here may be some problem
		// thinking about using goroutine pool
		go builder.report(ctx, inv, dur)
		return res
	}
}

func (builder *FilterChainBuilder) report(ctx context.Context, inv *orm.Invocation, dur time.Duration) {
	// start a transaction, we don't record it
	if strings.HasPrefix(inv.Method, "Begin") {
		return
	}
	if inv.Method == "Commit" || inv.Method == "Rollback" {
		builder.reportTxn(ctx, inv)
		return
	}
	builder.summaryVec.WithLabelValues(inv.Method, inv.GetTableName(), strconv.Itoa(int(dur)),
		strconv.FormatBool(inv.InsideTx), inv.TxName)
}

func (builder *FilterChainBuilder) reportTxn(ctx context.Context, inv *orm.Invocation) {
	dur := time.Now().Sub(inv.TxStartTime) / time.Millisecond
	builder.summaryVec.WithLabelValues(inv.Method, inv.TxName, strconv.Itoa(int(dur)),
		strconv.FormatBool(inv.InsideTx), inv.TxName)
}
