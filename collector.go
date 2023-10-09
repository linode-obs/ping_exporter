package main

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
	namespace = "ping_"
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

func pingHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	target := params.Get("target")
	timeout := time.Second * 10
	interval := time.Second
	count := 5
	size := 56
	ttl := 64
	proto := "ip4:icmp"
	for k, v := range params {
		k = strings.ToLower(k)
		if k == "target" {
			target = v[0]
		}
		if k == "timeout" {
			if strings.HasSuffix(v[0], "s") {
				value, err := time.ParseDuration(v[0])
				if err != nil {
					log.Error(err)
				}
				timeout = value
			} else {
				log.Errorf("expecting time duration in seconds example: 5s - got: %v", v[0])
			}
		}
		if k == "interval" {
			if strings.HasSuffix(v[0], "s") {
				value, err := time.ParseDuration(v[0])
				if err != nil {
					log.Error(err)
				}
				interval = value
			} else {
				log.Warnf("expecting time duration in seconds example: 5s - got: %v - using 1s default", v[0])
			}
		}
		if k == "count" {
			value, err := strconv.Atoi(v[0])
			if err != nil {
				log.Error(err)
			}
			if value > 0 {
				count = value
			}
		}
		if k == "size" {
			value, err := strconv.Atoi(v[0])
			if err != nil {
				log.Error(err)
			}
			if value < 1024 {
				size = value
			}
		}
		if k == "ttl" {
			value, err := strconv.Atoi(v[0])
			if err != nil {
				log.Error(err)
			}
			ttl = value
		}
		if k == "proto" {
			value := strings.ToLower(v[0])
			proto = value
		}

	}

	start := time.Now()
	registry := prometheus.NewRegistry()
	registry.MustRegister(probeDurationGauge)
	registry.MustRegister(minGauge)
	registry.MustRegister(maxGauge)
	registry.MustRegister(avgGauge)
	registry.MustRegister(stddevGauge)
	registry.MustRegister(lossGauge)

	ra, err := net.ResolveIPAddr(proto, target)
	if err != nil {
		log.Error(err)
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
		return
	}

	pinger, err := probing.NewPinger(ra.IP.String())
	if err != nil {
		log.Error(err)
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
		return
	}

	pinger.Count = count
	pinger.Size = size
	pinger.Interval = interval
	pinger.Timeout = timeout
	pinger.TTL = ttl
	pinger.SetPrivileged(false)

	err = pinger.Run()
	if err != nil {
		log.Error(err)
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
		return
	} else {
		stats := pinger.Statistics()

		minGauge.Set(stats.MinRtt.Seconds())
		avgGauge.Set(stats.AvgRtt.Seconds())
		maxGauge.Set(stats.MaxRtt.Seconds())
		stddevGauge.Set(float64(stats.StdDevRtt))
		lossGauge.Set(stats.PacketLoss)

		duration := time.Since(start).Seconds()
		probeDurationGauge.Set(duration)
	}
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
