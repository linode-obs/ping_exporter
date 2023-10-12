# Prometheus Ping Exporter

Yet another ping exporter.

Features:

* Uses Prometheus [multi-target export pattern](https://prometheus.io/docs/guides/multi-target-exporter/)
* IPv4 and IPv6
* UDP and ICMP

This ping exporter is simple and fills a gap in the ping exporter market just by supporting the multi-target export pattern. No additional prober configuration is needed, as many other ping exporters require. A scrape job will take care of all configuration.

- [Prometheus Ping Exporter](#prometheus-ping-exporter)
  - [Parameters](#parameters)
  - [Metrics](#metrics)
  - [Other Ping Exporters](#other-ping-exporters)
  - [Contributors](#contributors)

## Parameters

| Parameter Name     | Description                                                                                                                               | Default | Acceptable Values                                         |
| ------------------ | ----------------------------------------------------------------------------------------------------------------------------------------- | ------- | --------------------------------------------------------- |
| `target`           | What to ping                                                                                                                              | none    | Any hostname or IPv4/v6 address                           |
| `timeout`          | How long the entire ping job should run before returning                                                                                  | 10s     | Any `time.Duration` value                                 |
| `interval`         | How long to wait between pings                                                                                                            | 1s      | Any `time.Duration` value                                 |
| `count`            | How many pings to send                                                                                                                    | 5       | Any integer value                                         |
| `size`             | The size of the packet                                                                                                                    | 56      | Any integer value between 24 and 1024                     |
| `TTL`              | TTL of the packet                                                                                                                         | 64      | Any `time.Duration` value                                 |
| `protocol`, `prot` | IPv4 or IPv6                                                                                                                              | 1s      | `v6`, `6`, `ip6` (all other values considered to be IPv4) |
| `packet`           | UDP or ICMP (ICMP [requires root](https://pkg.go.dev/github.com/prometheus-community/pro-bing@v0.3.0#Pinger.SetPrivileged) in most cases) | `icmp`  | `icmp` (all other values considered to be `udp`)          |

## Metrics

| Metric Name            | Type  | Description                                            |
| ---------------------- | ----- | ------------------------------------------------------ |
| ping_duration_seconds  | gauge | Returns how long the probe took to complete in seconds |
| ping_loss_ratio        | gauge | Packet loss from 0 to 100                              |
| ping_rtt_avg_seconds   | gauge | Mean round trip time                                   |
| ping_rtt_max_seconds   | gauge | Worst round trip time                                  |
| ping_rtt_min_seconds   | gauge | Best round trip time                                   |
| ping_rtt_std_deviation | gauge | Standard deviation                                     |
| ping_success           | gauge | Returns whether the ping succeeded                     |

## Other Ping Exporters

There are many other good ping exporters, please give them a look too:

* [ping_exporter by czerwonk](https://github.com/czerwonk/ping_exporter)
* [mtr-exporter by mgumz](https://github.com/mgumz/mtr-exporter)
* [smokeping_prober by SuperQ](https://github.com/SuperQ/smokeping_prober)
* [ping-exporter by knsd](https://github.com/knsd/ping-exporter)

## Contributors

Special thanks to [stalloneclone](https://github.com/stalloneclone) for supplying the idea and initial code for the project!
