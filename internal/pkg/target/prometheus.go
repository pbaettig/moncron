package target

import (
	"fmt"

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

func (p PrometheusPushgateway) Name() string {
	return "prometheus-pushgateway"
}

func (p PrometheusPushgateway) Push(r *run.Command) error {
	if r == nil {
		return fmt.Errorf("nothing to push")
	}

	p.exitCodeGauge.WithLabelValues(r.Name).Set(float64(r.Result.ExitCode))
	p.durationGauge.WithLabelValues(r.Name).Set(r.Result.WallTime.Seconds())
	p.maxRssGauge.WithLabelValues(r.Name).Set(float64(r.Result.MaxRssBytes))
	p.lastExecutionGauge.WithLabelValues(r.Name).SetToCurrentTime()
	p.userTimeGauge.WithLabelValues(r.Name).Set(float64(r.Result.UserTime.Nano()) / 1000 / 1000 / 1000)
	p.systemTimeGauge.WithLabelValues(r.Name).Set(float64(r.Result.SystemTime.Nano()) / 1000 / 1000 / 1000)

	return push.New(p.URL, r.Name).Gatherer(p.Registry).Push()
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
