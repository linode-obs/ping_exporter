package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wbollock/ping_exporter/collector"
	log "github.com/sirupsen/logrus"
)

const (
	version = "0.1.0"
)

var (
	listenAddress = flag.String("web.listen-address", ":9141", "Address to listen on for telemetry")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics")
	showVersion   = flag.Bool("version", false, "show version information")
)

func printVersion() {
	fmt.Println("ping_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("multi-target ICMP prometheus exporter")
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	log.Info("Listening on ", *listenAddress)
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc(("/probe"), func(w http.ResponseWriter, r *http.Request) {
		collector.PingHandler(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html>
			<head><title>Ping Exporter</title></head>
			<body>
			<h1>Ping Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
		if err != nil {
			log.Fatal(err)
		}
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
