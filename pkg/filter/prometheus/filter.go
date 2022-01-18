package prometheus

import (
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/bhojpur/web/pkg/context"
	bhojpur "github.com/bhojpur/web/pkg/engine"
	web "github.com/bhojpur/web/pkg/engine"
)

// FilterChainBuilder is an extension point,
// when we want to support some configuration,
// please use this structure
type FilterChainBuilder struct {
}

// FilterChain returns a FilterFunc. The filter will records some metrics
func (builder *FilterChainBuilder) FilterChain(next web.FilterFunc) web.FilterFunc {
	summaryVec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:      "bhojpur",
		Subsystem: "http_request",
		ConstLabels: map[string]string{
			"server":  web.BasConfig.ServerName,
			"env":     web.BasConfig.RunMode,
			"appname": web.BasConfig.AppName,
		},
		Help: "The statistics info for HTTP request",
	}, []string{"pattern", "method", "status", "duration"})

	prometheus.MustRegister(summaryVec)

	registerBuildInfo()

	return func(ctx *context.Context) {
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
			"appname":        web.BasConfig.AppName,
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

func report(dur time.Duration, ctx *context.Context, vec *prometheus.SummaryVec) {
	status := ctx.Output.Status
	ptn := ctx.Input.GetData("RouterPattern").(string)
	ms := dur / time.Millisecond
	vec.WithLabelValues(ptn, ctx.Input.Method(), strconv.Itoa(status), strconv.Itoa(int(ms))).Observe(float64(ms))
}
