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
    - [/probe](#probe)
    - [/metrics](#metrics-1)
  - [Example Scrape Job](#example-scrape-job)
  - [Installation](#installation)
    - [Debian/RPM package](#debianrpm-package)
    - [Docker](#docker)
    - [Binary](#binary)
    - [Source](#source)
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

### /probe

| Metric Name            | Type  | Description                                            |
| ---------------------- | ----- | ------------------------------------------------------ |
| ping_duration_seconds  | gauge | Returns how long the probe took to complete in seconds |
| ping_loss_ratio        | gauge | Packet loss from 0 to 100                              |
| ping_rtt_avg_seconds   | gauge | Mean round trip time                                   |
| ping_rtt_max_seconds   | gauge | Worst round trip time                                  |
| ping_rtt_min_seconds   | gauge | Best round trip time                                   |
| ping_rtt_std_deviation | gauge | Standard deviation                                     |
| ping_success           | gauge | Returns whether the ping succeeded                     |

### /metrics

Standard Prometheus webserver metrics, plus `ping_exporter_version_info`.

## Example Scrape Job

```yaml
scrape_configs:
  - job_name: 'ping_exporter_metrics'
    scrape_interval: 15s
    scrape_timeout: 15s
    static_configs:
      - targets: ['0.0.0.0:9141']

  - job_name: 'prober'
    metrics_path: '/probe'
    scrape_interval: 15s
    scrape_timeout: 15s
    params:
      protocol: ['4']
      count: ['5']
      interval: ['1s']
    honor_labels: True
    static_configs:
    - targets:
      - google.com
      - linode.com
    relabel_configs:
    # Set the exporter's target
    - source_labels: [__address__]
      target_label: __param_target
    # Set address label to instance
    - source_labels: [__address__]
      target_label: instance
    # Actually talk to the blackbox exporter
    - target_label: __address__
      replacement: 0.0.0.0:9141
    # If we set a custom instance label, write it to the
    # expected instance label
    - source_labels: [__instance]
      target_label: instance
      regex: '(.+)'
      replacement: '${1}'
```

## Installation

### Debian/RPM package

Substitute `{{ version }}` for your desired release.

```bash
wget https://github.com/wbollock/ping_exporter/releases/download/v{{ version }}/prometheus-ping-exporter_{{ version }}_linux_amd64.{deb,rpm}
{dpkg,rpm} -i prometheus-ping-exporter_{{ version }}_linux_amd64.{deb,rpm}
```

### Docker

```console
sudo docker run \
--privileged \
ghcr.io/wbollock/ping_exporter
```

### Binary

```bash
wget https://github.com/wbollock/ping_exporter/releases/download/v{{ version }}/ping_exporter_{{ version }}_Linux_x86_64.tar.gz
tar xvf ping_exporter_{{ version }}_Linux_x86_64.tar.gz
./ping_exporter/prometheus-ping-exporter
```

### Source

```bash
wget https://github.com/wbollock/ping_exporter/archive/refs/tags/v{{ version }}.tar.gz
tar xvf ping_exporter-{{ version }}.tar.gz
cd ./ping_exporter-{{ version }}
go build ping_exporter.go
./ping_exporter.go
```

## Other Ping Exporters

There are many other good ping exporters, please give them a look too:

* [ping_exporter by czerwonk](https://github.com/czerwonk/ping_exporter)
* [mtr-exporter by mgumz](https://github.com/mgumz/mtr-exporter)
* [smokeping_prober by SuperQ](https://github.com/SuperQ/smokeping_prober)
* [ping-exporter by knsd](https://github.com/knsd/ping-exporter)

## Contributors

Special thanks to [stalloneclone](https://github.com/stalloneclone) for supplying the idea and initial code for the project!
