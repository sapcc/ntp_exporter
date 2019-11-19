# ntp\_exporter

This is a Prometheus exporter that, when running on a node, checks the drift
of that node's clock against a given NTP server.

## Installation

Compile `make && make install` or `docker build`. The binary can also be
installed with `go get`:

```bash
go get github.com/sapcc/ntp_exporter
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

- `server`: NTP server to use
- `protocol`: NTP protocol version (2, 3 or 4)
- `duration`: duration of measurements in case of high drift

For example:

```sh
$ curl 'http://localhost:9559/metrics?server=ntp.example.com&protocol=4&duration=10s'
```
