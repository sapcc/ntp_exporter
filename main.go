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
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var version string // will be substituted at compile-time

func main() {
	var (
		showVersion            = flag.Bool("version", false, "Print version information.")
		listenAddress          = flag.String("web.listen-address", ":9559", "Address on which to expose metrics and web interface.")
		metricsPath            = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		ntpServer              = flag.String("ntp.server", "", "NTP server to use (required).")
		ntpProtocolVersion     = flag.Int("ntp.protocol-version", 4, "NTP protocol version to use.")
		ntpMeasurementDuration = flag.Duration("ntp.measurement-duration", 30*time.Second, "Duration of measurements in case of high (>10ms) drift.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *ntpServer == "" {
		log.Fatalln("no NTP server specified, see -ntp.server")
	}
	if *ntpProtocolVersion < 2 || *ntpProtocolVersion > 4 {
		log.Fatalf("invalid NTP protocol version %d; must be 2, 3, or 4", *ntpProtocolVersion)
	}

	log.Infoln("starting ntp_exporter", version)
	prometheus.MustRegister(Collector{*ntpServer, *ntpProtocolVersion, *ntpMeasurementDuration})
	handler := promhttp.HandlerFor(prometheus.DefaultGatherer,
		promhttp.HandlerOpts{ErrorLog: log.NewErrorLogger()})

	http.Handle(*metricsPath, prometheus.InstrumentHandler("prometheus", handler))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>NTP Exporter</title></head>
			<body>
			<h1>NTP Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Infoln("listening on", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
