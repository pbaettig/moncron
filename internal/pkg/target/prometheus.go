package target

import (
	"fmt"

	"github.com/pbaettig/moncron/internal/pkg/model"
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

func (p PrometheusPushgateway) Push(r model.JobRun) error {

	p.exitCodeGauge.WithLabelValues(r.Name, r.Host.Name).Set(float64(r.Result.ExitCode))
	p.durationGauge.WithLabelValues(r.Name, r.Host.Name).Set(r.Result.WallTime.Seconds())
	p.maxRssGauge.WithLabelValues(r.Name, r.Host.Name).Set(float64(r.Result.MaxRssBytes))
	p.lastExecutionGauge.WithLabelValues(r.Name, r.Host.Name).SetToCurrentTime()
	p.userTimeGauge.WithLabelValues(r.Name, r.Host.Name).Set(float64(r.Result.UserTime.Seconds()))
	p.systemTimeGauge.WithLabelValues(r.Name, r.Host.Name).Set(float64(r.Result.SystemTime.Seconds()))

	return push.New(p.URL, fmt.Sprintf("%s/%s", r.Name, r.Host.Name)).Gatherer(p.Registry).Push()
}

func NewPrometheusPushgateway(url string) PrometheusPushgateway {
	p := PrometheusPushgateway{URL: url, Registry: prometheus.NewRegistry()}

	p.exitCodeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "exit_code",
	}, []string{"job_name", "host_name"})

	p.durationGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "duration_seconds",
	}, []string{"job_name", "host_name"})

	p.maxRssGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "max_rss_bytes",
	}, []string{"job_name", "host_name"})

	p.lastExecutionGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "last_execution",
	}, []string{"job_name", "host_name"})

	p.userTimeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "user_time_seconds",
	}, []string{"job_name", "host_name"})

	p.systemTimeGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "system_time_seconds",
	}, []string{"job_name", "host_name"})

	p.Registry.MustRegister(p.exitCodeGauge)
	p.Registry.MustRegister(p.durationGauge)
	p.Registry.MustRegister(p.maxRssGauge)
	p.Registry.MustRegister(p.lastExecutionGauge)
	p.Registry.MustRegister(p.userTimeGauge)
	p.Registry.MustRegister(p.systemTimeGauge)

	return p
}
