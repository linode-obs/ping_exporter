package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type PingMetrics struct {
	PingSuccessGauge   prometheus.Gauge
	PingTimeoutGauge   prometheus.Gauge
	ProbeDurationGauge prometheus.Gauge
	MinGauge           prometheus.Gauge
	MaxGauge           prometheus.Gauge
	AvgGauge           prometheus.Gauge
	StddevGauge        prometheus.Gauge
	LossGauge          prometheus.Gauge
}
