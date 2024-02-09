package target

import (
	"github.com/pbaettig/moncron/internal/pkg/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
)

var (
	namespace = "moncron"
	subsystem = "job"
)

type PrometheusPushgateway struct {
	ResultTarget
	// JobName            string
	URL                string
	Registry           *prometheus.Registry
	exitCodeGauge      *prometheus.GaugeVec
	durationGauge      *prometheus.GaugeVec
	maxRssGauge        *prometheus.GaugeVec
	lastExecutionGauge *prometheus.GaugeVec
	userTimeGauge      *prometheus.GaugeVec
	systemTimeGauge    *prometheus.GaugeVec
}

func (p PrometheusPushgateway) Push(jobName string, r run.CommandResult) error {
	p.exitCodeGauge.WithLabelValues(jobName).Set(float64(r.ExitCode))
	p.durationGauge.WithLabelValues(jobName).Set(r.WallTime.Seconds())
	p.maxRssGauge.WithLabelValues(jobName).Set(float64(r.MaxRSS))
	p.lastExecutionGauge.WithLabelValues(jobName).SetToCurrentTime()
	p.userTimeGauge.WithLabelValues(jobName).Set(float64(r.UserTimeMilli) / 1000)
	p.systemTimeGauge.WithLabelValues(jobName).Set(float64(r.SystemTimeMilli) / 1000)

	return push.New(p.URL, jobName).Gatherer(p.Registry).Push()
}

func NewPrometheusPushgateway(url string) PrometheusPushgateway {
	p := PrometheusPushgateway{URL: url, Registry: prometheus.NewRegistry()}

	p.exitCodeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "exit_code",
	}, []string{"name"})

	p.durationGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "duration_seconds",
	}, []string{"name"})

	p.maxRssGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "max_rss_bytes",
	}, []string{"name"})

	p.lastExecutionGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "last_execution",
	}, []string{"name"})

	p.userTimeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "user_time_seconds",
	}, []string{"name"})

	p.systemTimeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "system_time_seconds",
	}, []string{"name"})

	p.Registry.MustRegister(p.exitCodeGauge)
	p.Registry.MustRegister(p.durationGauge)
	p.Registry.MustRegister(p.maxRssGauge)
	p.Registry.MustRegister(p.lastExecutionGauge)
	p.Registry.MustRegister(p.userTimeGauge)
	p.Registry.MustRegister(p.systemTimeGauge)

	return p
}
