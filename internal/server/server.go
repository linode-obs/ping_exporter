package server

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/wbollock/ping_exporter/internal/collector"
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

func SetupServer(metricsPath string) http.Handler {
	mux := http.NewServeMux()

	mux.Handle(metricsPath, promhttp.Handler())
	mux.HandleFunc("/probe", collector.PingHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf(defaultHTML, metricsPath)
		_, err := w.Write([]byte(response))
		if err != nil {
			log.WithError(err).Error("Failed to write main page response")
		}
	})

	return mux
}
