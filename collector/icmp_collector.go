package collector

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	namespace       = "ping_"
	defaultTimeout  = time.Second * 10
	defaultInterval = time.Second
	defaultCount    = 5
	defaultSize     = 56
	defaultTTL      = 64
	protoDefault    = "ip4:icmp"
	maxPacketSize   = 1024
)

var (
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

type pingParams struct {
	target   string
	timeout  time.Duration
	interval time.Duration
	count    int
	size     int
	ttl      int
	proto    string
}

func parseParams(r *http.Request) pingParams {
	params := r.URL.Query()
	p := pingParams{
		target:   params.Get("target"),
		timeout:  defaultTimeout,
		interval: defaultInterval,
		count:    defaultCount,
		size:     defaultSize,
		ttl:      defaultTTL,
		proto:    protoDefault,
	}

	for k, v := range params {
		switch strings.ToLower(k) {
		case "target":
			p.target = v[0]
		case "timeout":
			if duration, err := time.ParseDuration(v[0]); err == nil {
				p.timeout = duration
			} else {
				log.Errorf("Expected duration in seconds (e.g., 5s). Got: %v", v[0])
			}
		case "interval":
			if duration, err := time.ParseDuration(v[0]); err == nil {
				p.interval = duration
			} else {
				log.Warnf("Expected duration in seconds (e.g., 5s). Got: %v. Using default 1s.", v[0])
			}
		case "count":
			if count, err := strconv.Atoi(v[0]); err == nil && count > 0 {
				p.count = count
			}
		case "size":
			if size, err := strconv.Atoi(v[0]); err == nil && size < maxPacketSize {
				p.size = size
			}
		case "ttl":
			if ttl, err := strconv.Atoi(v[0]); err == nil {
				p.ttl = ttl
			}
		case "proto":
			p.proto = strings.ToLower(v[0])
		}
	}

	return p
}

func serveMetricsWithError(w http.ResponseWriter, r *http.Request, registry *prometheus.Registry) {
	if h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{}); h != nil {
		h.ServeHTTP(w, r)
	}
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	p := parseParams(r)
	start := time.Now()

	registry := prometheus.NewRegistry()
	registry.MustRegister(probeDurationGauge, minGauge, maxGauge, avgGauge, stddevGauge, lossGauge)

	// TODO: ensure ResolveIPAddr is the best way to do lookups
	ra, err := net.ResolveIPAddr(p.proto, p.target)
	if err != nil {
		log.Error(err)
		serveMetricsWithError(w, r, registry)
		return
	}

	pinger, err := probing.NewPinger(ra.IP.String())
	if err != nil {
		log.Error(err)
		serveMetricsWithError(w, r, registry)
		return
	}

	pinger.Count = p.count
	pinger.Size = p.size
	pinger.Interval = p.interval
	pinger.Timeout = p.timeout
	pinger.TTL = p.ttl
	pinger.SetPrivileged(false)

	if err := pinger.Run(); err != nil {
		log.Error(err)
		serveMetricsWithError(w, r, registry)
		return
	}

	stats := pinger.Statistics()
	minGauge.Set(stats.MinRtt.Seconds())
	avgGauge.Set(stats.AvgRtt.Seconds())
	maxGauge.Set(stats.MaxRtt.Seconds())
	stddevGauge.Set(float64(stats.StdDevRtt))
	lossGauge.Set(stats.PacketLoss)
	probeDurationGauge.Set(time.Since(start).Seconds())

	serveMetricsWithError(w, r, registry)
}
