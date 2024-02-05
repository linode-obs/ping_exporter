package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/linode-obs/ping_exporter/internal/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	defaultLogLevel      = "info"
	defaultListenAddress = "0.0.0.0:9141"
	defaultMetricsPath   = "/metrics"
)

var (
	listenAddress = flag.String("web.listen-address", defaultListenAddress, "Address to listen on for telemetry")
	showVersion   = flag.Bool("version", false, "show version information")
	logLevel      = flag.String("log.level", defaultLogLevel,
		"Minimum Log level [debug, info]")

	// Build info for ping exporter itself, will be populated by linker during build
	Version   string
	BuildDate string
	Commit    string

	versionInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ping_exporter_version_info",
			Help: "Ping Exporter build information",
		},
		[]string{"version", "commit", "builddate"},
	)
)

func printVersion() {
	fmt.Printf("ping_exporter\n")
	fmt.Printf("Version:   %s\n", Version)
	fmt.Printf("BuildDate: %s\n", BuildDate)
	fmt.Printf("Commit:    %s\n", Commit)
	fmt.Printf("multi-target ICMP prometheus exporter\n")
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	versionInfo.WithLabelValues(Version, Commit, BuildDate).Set(1)
	prometheus.MustRegister(versionInfo)

	switch *logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.Debug("Log level set to debug")
	default:
		log.SetLevel(log.InfoLevel)
	}

	http.Handle(defaultMetricsPath, promhttp.Handler())
	http.Handle("/", server.SetupServer())

	log.Infof("Starting server on %s", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.WithError(err).Fatal("Failed to start the server")
	}
}
