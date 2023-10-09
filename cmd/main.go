package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/wbollock/ping_exporter/internal/collector"
)

const (
	version     = "0.1.0"
	defaultHTML = `<html>
			<head><title>Ping Exporter</title></head>
			<body>
			<h1>Ping Exporter</h1>
			<p><a href='%s'>Metrics</a></p>
			</body>
			</html>`
	defaultListenAddress = ":9141"
	defaultMetricsPath   = "/metrics"
	defaultLogLevel      = "info"
)

var (
	listenAddress = flag.String("web.listen-address", defaultListenAddress, "Address to listen on for telemetry")
	metricsPath   = flag.String("web.telemetry-path", defaultMetricsPath, "Path under which to expose metrics")
	showVersion   = flag.Bool("version", false, "show version information")
	logLevel      = flag.String("log.level", defaultLogLevel,
		"Minimum Log level [debug, info]")
)

func printVersion() {
	fmt.Printf("ping_exporter\nVersion: %s\nmulti-target ICMP prometheus exporter\n", version)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	response := fmt.Sprintf(defaultHTML, *metricsPath)
	_, err := w.Write([]byte(response))
	if err != nil {
		log.WithError(err).Error("Failed to write main page response")
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	switch *logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.Debug("Log level set to debug")
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.Info("Listening on address ", *listenAddress)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/probe", collector.PingHandler)
	http.HandleFunc("/", mainPage)

	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.WithError(err).Fatal("Failed to start the server")
	}
}
