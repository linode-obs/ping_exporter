package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/linode-obs/ping_exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func SetupServer() http.Handler {

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

	mux := http.NewServeMux()

	mux.Handle(defaultMetricsPath, promhttp.Handler())

	pingHandler := collector.PingHandler()

	mux.HandleFunc("/probe", pingHandler)

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
