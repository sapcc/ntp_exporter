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
	"fmt"
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	drift = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "ntp",
		Name:      "drift_seconds",
		Help:      "Difference between system time and NTP time.",
	})
	stratum = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "ntp",
		Name:      "stratum",
		Help:      "Stratum of NTP server.",
	})
	scrapeDuration = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "ntp",
		Name:      "scrape_duration_seconds",
		Help:      "ntp_exporter: Duration of a scrape job.",
	})
)

//Collector implements the prometheus.Collector interface.
type Collector struct {
	NtpServer          string
	NtpProtocolVersion int
}

//Describe implements the prometheus.Collector interface.
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	drift.Describe(ch)
	stratum.Describe(ch)
	scrapeDuration.Describe(ch)
}

//Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	err := c.measure()
	//only report data when measurement was successful
	if err == nil {
		drift.Collect(ch)
		stratum.Collect(ch)
		scrapeDuration.Collect(ch)
	} else {
		log.Errorln(err)
		return
	}
}

func (c Collector) measure() error {
	begin := time.Now()

	resp, err := ntp.Query(c.NtpServer, c.NtpProtocolVersion)
	if err != nil {
		return fmt.Errorf("couldn't get NTP drift: %s", err)
	}
	drift.Set(resp.ClockOffset.Seconds())
	stratum.Set(float64(resp.Stratum))

	scrapeDuration.Observe(time.Since(begin).Seconds())
	return nil
}
