# ntp\_exporter

This is a Prometheus exporter that, when running on a node, checks the drift
of that node's clock against a given NTP server.

## Installation

To build the binary:

```bash
make
```

The binary can also be installed with `go get`:

```bash
go get github.com/sapcc/ntp_exporter
```

To build the Docker container:

```bash
make GOLDFLAGS="-w -linkmode external -extldflags -static" && docker build .
```

## Usage

Command-line options: (Only `-ntp.server` or `-ntp.source` is required.)

```plain
Usage of ntp_exporter:
  -ntp.measurement-duration duration
     Duration of measurements in case of high (>10ms) drift. (default 30s)
  -ntp.protocol-version int
     NTP protocol version to use. (default 4)
  -ntp.server string
     NTP server to use (required).
  -ntp.source string
     source of information about ntp server (cli / http). (default "cli")
  -version
     Print version information.
  -web.listen-address string
     Address on which to expose metrics and web interface. (default ":9559")
  -web.telemetry-path string
     Path under which to expose metrics. (default "/metrics")
```
