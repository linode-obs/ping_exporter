package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/wbollock/ping_exporter/internal/collector"
	"github.com/wbollock/ping_exporter/internal/metrics"
)

const (
	defaultHTML = `<html>
			<head><title>Ping Exporter</title></head>
			<body>
			<h1>Ping Exporter</h1>
			<p><a href='%s'>Metrics</a></p>
			</body>
			</html>`
	defaultMetricsPath = "/metrics"
)

const namespace = "ping_"

var (
	pingSuccessGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "success",
		Help: "Returns whether the ping succeeded",
	})
	pingTimeoutGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "timeout",
		Help: "Returns whether the ping failed by timeout",
	})
	probeDurationGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "duration_seconds",
		Help: "Returns how long the probe took to complete in seconds",
	})
	minGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "rtt_min_seconds",
		Help: "Best round trip time",
	})
	maxGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "rtt_max_seconds",
		Help: "Worst round trip time",
	})
	avgGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "rtt_avg_seconds",
		Help: "Mean round trip time",
	})
	stddevGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "rtt_std_deviation",
		Help: "Standard deviation",
	})
	lossGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: namespace + "loss_ratio",
		Help: "Packet loss from 0 to 100",
	})
)

func SetupServer() http.Handler {
	mux := http.NewServeMux()

	mux.Handle(defaultMetricsPath, promhttp.Handler())

	var mutex sync.Mutex
	registry := prometheus.NewRegistry()

	metrics := metrics.PingMetrics{
		PingSuccessGauge:   pingSuccessGauge,
		PingTimeoutGauge:   pingTimeoutGauge,
		ProbeDurationGauge: probeDurationGauge,
		MinGauge:           minGauge,
		MaxGauge:           maxGauge,
		AvgGauge:           avgGauge,
		StddevGauge:        stddevGauge,
		LossGauge:          lossGauge,
	}

	registry.MustRegister(metrics.PingSuccessGauge, metrics.PingTimeoutGauge, metrics.ProbeDurationGauge, metrics.MinGauge, metrics.MaxGauge, metrics.AvgGauge, metrics.StddevGauge, metrics.LossGauge)
	mux.HandleFunc("/probe", collector.PingHandler(registry, metrics, &mutex))

	// for non-standard web servers, need to register handlers
	mux.HandleFunc("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.HandleFunc("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.HandleFunc("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.HandleFunc("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.HandleFunc("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf(defaultHTML, defaultMetricsPath)
		_, err := w.Write([]byte(response))
		if err != nil {
			log.WithError(err).Error("Failed to write main page response")
		}
	})

	return mux
}
