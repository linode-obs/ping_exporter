package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/wbollock/ping_exporter/internal/server"
)

const (
	version              = "0.1.0"
	defaultLogLevel      = "info"
	defaultListenAddress = ":9141"
	defaultMetricsPath   = "/metrics"
)

var (
	listenAddress = flag.String("web.listen-address", defaultListenAddress, "Address to listen on for telemetry")
	showVersion   = flag.Bool("version", false, "show version information")
	logLevel      = flag.String("log.level", defaultLogLevel,
		"Minimum Log level [debug, info]")
)

func printVersion() {
	fmt.Printf("ping_exporter\nVersion: %s\nmulti-target ICMP prometheus exporter\n", version)
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

	handler := server.SetupServer()

	log.Infof("Starting server on %s", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, handler); err != nil {
		log.WithError(err).Fatal("Failed to start the server")
	}
}
