/*******************************************************************************
*
* Copyright 2017 SAP SE
* Copyright 2015 The Prometheus Authors
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You should have received a copy of the License along with this
* program. If not, you may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
*******************************************************************************/

package main

import (
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
)

var (
	showVersion            bool
	listenAddress          string
	metricsPath            string
	ntpServer              string
	ntpProtocolVersion     int
	ntpMeasurementDuration time.Duration
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
			log.Fatalf("invalid NTP protocol version %d; must be 2, 3, or 4", ntpProtocolVersion)
		}
	}

	logger.Println("starting ntp_exporter", version)

	http.Handle(metricsPath, http.HandlerFunc(handlerMetrics))
	http.HandleFunc("/", handlerDefault)

	logger.Println("listening on", listenAddress)
	err := http.ListenAndServe(listenAddress, nil) //nolint: gosec // no timeout is required
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.BoolVar(&showVersion, "version", false, "Print version information.")
	flag.StringVar(&listenAddress, "web.listen-address", ":9559", "Address on which to expose metrics and web interface.")
	flag.StringVar(&metricsPath, "web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&ntpServer, "ntp.server", "", "NTP server to use (required).")
	flag.IntVar(&ntpProtocolVersion, "ntp.protocol-version", 4, "NTP protocol version to use.")
	flag.DurationVar(&ntpMeasurementDuration, "ntp.measurement-duration", 30*time.Second, "Duration of measurements in case of high (>10ms) drift.")
	flag.StringVar(&ntpSource, "ntp.source", "cli", "source of information about ntp server (cli / http).")
	flag.Parse()
}

func handlerMetrics(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	s := ntpServer
	p := ntpProtocolVersion
	d := ntpMeasurementDuration

	if ntpSource == "http" {
		for _, i := range []string{"target", "protocol", "duration"} {
			if r.URL.Query().Get(i) == "" {
				http.Error(w, fmt.Sprintf("Get parameter is empty: %s", i), http.StatusBadRequest)
				return
			}
		}

		s = r.URL.Query().Get("target")

		if v, err := strconv.ParseInt(r.URL.Query().Get("protocol"), 10, 32); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if v < 2 || v > 4 {
			http.Error(w, fmt.Sprintf("invalid NTP protocol version %d; must be 2, 3, or 4", v), http.StatusBadRequest)
			return
		} else {
			p = int(v)
		}

		if t, err := time.ParseDuration(r.URL.Query().Get("duration")); err == nil {
			d = t
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(CollectorInitial(s, p, d))
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
