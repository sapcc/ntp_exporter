// SPDX-FileCopyrightText: 2017 SAP SE or an SAP affiliate company
// SPDX-FileCopyrightText: 2015 The Prometheus Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sapcc/go-api-declarations/bininfo"
	"github.com/sapcc/go-bits/httpext"
	"github.com/sapcc/go-bits/must"
	_ "go.uber.org/automaxprocs"
)

var (
	showVersion            bool
	listenAddress          string
	metricsPath            string
	ntpServer              string
	ntpProtocolVersion     int
	ntpMeasurementDuration time.Duration
	ntpHighDrift           time.Duration
	ntpSource              string
)

var logger = log.New(os.Stderr, "", log.LstdFlags)

var version = bininfo.VersionOr("unknown")
var buildDate = bininfo.BuildDateOr("unknown")
var revision = bininfo.CommitOr("unknown")

func main() {
	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if ntpSource == "cli" {
		if ntpServer == "" {
			logger.Fatalln("no NTP server specified, see -ntp.server")
		}

		if ntpProtocolVersion < 2 || ntpProtocolVersion > 4 {
			logger.Fatalf("invalid NTP protocol version %d; must be 2, 3, or 4", ntpProtocolVersion)
		}
	}

	logger.Println("starting ntp_exporter", version)

	mux := http.NewServeMux()
	mux.Handle(metricsPath, http.HandlerFunc(handlerMetrics))
	mux.HandleFunc("/", handlerDefault)
	ctx := httpext.ContextWithSIGINT(context.Background(), 1*time.Second)
	must.Succeed(httpext.ListenAndServeContext(ctx, listenAddress, mux))
}

func init() {
	flag.BoolVar(&showVersion, "version", false, "Print version information.")
	flag.StringVar(&listenAddress, "web.listen-address", ":9559", "Address on which to expose metrics and web interface.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&ntpServer, "ntp.server", "", "NTP server to use (required).")
	flag.IntVar(&ntpProtocolVersion, "ntp.protocol-version", 4, "NTP protocol version to use.")
	flag.DurationVar(&ntpMeasurementDuration, "ntp.measurement-duration", 30*time.Second, "Duration of measurements in case of high drift.")
	flag.DurationVar(&ntpHighDrift, "ntp.high-drift", 10*time.Millisecond, "Absolute high drift threshold.")
	flag.StringVar(&ntpSource, "ntp.source", "cli", "source of information about ntp server (cli / http).")
	flag.Parse()
}

func handlerMetrics(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s := ntpServer
	p := ntpProtocolVersion
	d := ntpMeasurementDuration
	hd := ntpHighDrift

	if d < 0 {
		logger.Fatal("ntp.measurement-duration cannot be negative")
	}
	if hd < 0 {
		logger.Fatal("ntp.high-drift cannot be negative")
	}

	if ntpSource == "http" {
		query := r.URL.Query()
		for _, i := range []string{"target", "protocol", "duration"} {
			if query.Get(i) == "" {
				http.Error(w, "Get parameter is empty: "+i, http.StatusBadRequest)
				return
			}
		}

		s = query.Get("target")

		if v, err := strconv.ParseInt(query.Get("protocol"), 10, 32); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if v < 2 || v > 4 {
			http.Error(w, fmt.Sprintf("invalid NTP protocol version %d; must be 2, 3, or 4", v), http.StatusBadRequest)
			return
		} else {
			p = int(v)
		}

		if t, err := time.ParseDuration(query.Get("duration")); err == nil {
			if t < 0 {
				http.Error(w, "duration cannot be negative", http.StatusBadRequest)
				return
			}
			d = t
		} else {
			http.Error(w, "while parsing duration: "+err.Error(), http.StatusBadRequest)
			return
		}

		if query.Get("high_drift") != "" {
			if u, err := time.ParseDuration(query.Get("high_drift")); err == nil {
				if u < 0 {
					http.Error(w, "high_drift cannot be negative", http.StatusBadRequest)
					return
				}
				hd = u
			} else {
				http.Error(w, "while parsing high_drift: "+err.Error(), http.StatusBadRequest)
				return
			}
		}
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(CollectorInitial(s, p, d, hd))
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: logger})
	h.ServeHTTP(w, r)
}

func handlerDefault(w http.ResponseWriter, r *http.Request) {
	//nolint:errcheck
	w.Write([]byte(`<html>
			<head><title>NTP Exporter</title></head>
			<body>
			<h1>NTP Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			<p>Version: ` + version + `</p>
			<p>Revision: ` + revision + `</p>
			<p>Build date: ` + buildDate + `</p>
			</body>
			</html>`))
}
