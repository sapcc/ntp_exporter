# ntp\_exporter

This is a Prometheus exporter that, when running on a node, checks the drift
of that node's clock against a given NTP server.

## Installation

Compile `make && make install` or `docker build`. The binary can also be
installed with `go get`:

```bash
go get github.com/sapcc/ntp_exporter
```

We also publish pre-built images on Docker Hub as
[sapcc/ntp-exporter](https://hub.docker.com/r/sapcc/ntp-exporter):

```bash
docker pull sapcc/ntp-exporter:v1.1.3
```

## Usage

Common command-line options:

```
-ntp.source string
   source of information about ntp server (cli / http). (default "cli")
-version
   Print version information.
-web.listen-address string
   Address on which to expose metrics and web interface. (default ":9559")
-web.telemetry-path string
   Path under which to expose metrics. (default "/metrics")
```

### Mode 1: Fixed NTP server

By default, or when the option `-ntp.source cli` is specified, the NTP server
and connection options is defined by command-line options:

```
-ntp.measurement-duration duration
   Duration of measurements in case of high (>10ms) drift. (default 30s)
-ntp.protocol-version int
   NTP protocol version to use. (default 4)
-ntp.server string
   NTP server to use (required).
```

### Mode 2: Variable NTP server

When the option `-ntp.source http` is specified, the NTP server and connection
options are obtained from the query parameters on each `GET /metrics` HTTP
request:

- `target`: NTP server to use
- `protocol`: NTP protocol version (2, 3 or 4)
- `duration`: duration of measurements in case of high drift

For example:

```sh
$ curl 'http://localhost:9559/metrics?target=ntp.example.com&protocol=4&duration=10s'
```

## Frequently asked questions (FAQ)

### Is there a metric for checking that the exporter is working?

Several people have suggested adding a metric like `ntp_up` that's always 1, so
that people can alert on `absent(ntp_up)` or something like that. This is not
necessary. [Prometheus already generates such a metric by itself during
scraping.](https://prometheus.io/docs/concepts/jobs_instances/) A suitable
alert expression could look like

```
up{job="ntp_exporter",instance="example.com:9559"} != 1 or absent(up{job="ntp_exporter",instance="example.com:9559"})
```

but the concrete labels will vary depending on your setup and scrape configuration.
