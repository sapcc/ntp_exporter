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
	"log"
	"sort"
	"time"

	"github.com/beevik/ntp"
	"github.com/prometheus/client_golang/prometheus"
)

func CollectorInitial(target string, protocol int, duration time.Duration) Collector {
	return Collector{
		NtpServer:              target,
		NtpProtocolVersion:     protocol,
		NtpMeasurementDuration: duration,
		buildInfo: prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Namespace:   "ntp",
			Name:        "build_info",
			Help:        "ntp_exporter version.",
			ConstLabels: prometheus.Labels{"version": version, "revision": revision, "build_date": buildDate},
		}, func() float64 { return 1 }),
		drift: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "drift_seconds",
			Help:      "Difference between system time and NTP time.",
		}, []string{"server"}),
		stratum: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "stratum",
			Help:      "Stratum of NTP server.",
		}, []string{"server"}),
		rtt: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "rtt_seconds",
			Help:      "Round-trip time to NTP server.",
		}, []string{"server"}),
		referenceTime: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "reference_timestamp_seconds",
			Help:      "Reference time of NTP server (UNIX Timestamp).",
		}, []string{"server"}),
		rootDelay: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "root_delay_seconds",
			Help:      "Root Delay of NTP server.",
		}, []string{"server"}),
		rootDispersion: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "root_dispersion_seconds",
			Help:      "Root Dispersion of NTP server.",
		}, []string{"server"}),
		rootDistance: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "root_distance_seconds",
			Help:      "Distance to Root NTP server.",
		}, []string{"server"}),
		precision: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "precision_seconds",
			Help:      "Precision of NTP server.",
		}, []string{"server"}),
		leap: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "ntp",
			Name:      "leap",
			Help:      "Leap second indicator.",
		}, []string{"server"}),
		scrapeDuration: prometheus.NewSummary(prometheus.SummaryOpts{
			Namespace: "ntp",
			Name:      "scrape_duration_seconds",
			Help:      "ntp_exporter: Duration of a scrape job.",
		}),
	}
}

// Collector implements the prometheus.Collector interface.
type Collector struct {
	NtpServer              string
	NtpProtocolVersion     int
	NtpMeasurementDuration time.Duration
	buildInfo              prometheus.GaugeFunc
	stratum                *prometheus.GaugeVec
	drift                  *prometheus.GaugeVec
	rtt                    *prometheus.GaugeVec
	referenceTime          *prometheus.GaugeVec
	rootDelay              *prometheus.GaugeVec
	rootDispersion         *prometheus.GaugeVec
	rootDistance           *prometheus.GaugeVec
	precision              *prometheus.GaugeVec
	leap                   *prometheus.GaugeVec
	scrapeDuration         prometheus.Summary
}

// A single measurement returned by ntp server
type Measurement struct {
	clockOffset    float64
	stratum        float64
	rtt            float64
	referenceTime  float64
	rootDelay      float64
	rootDispersion float64
	rootDistance   float64
	precision      float64
	leap           float64
}

// Describe implements the prometheus.Collector interface.
func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	c.buildInfo.Describe(ch)
	c.drift.Describe(ch)
	c.stratum.Describe(ch)
	c.rtt.Describe(ch)
	c.referenceTime.Describe(ch)
	c.rootDelay.Describe(ch)
	c.rootDispersion.Describe(ch)
	c.rootDistance.Describe(ch)
	c.precision.Describe(ch)
	c.leap.Describe(ch)
	c.scrapeDuration.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (c Collector) Collect(ch chan<- prometheus.Metric) {
	err := c.measure()
	//only report data when measurement was successful
	if err == nil {
		c.buildInfo.Collect(ch)
		c.drift.Collect(ch)
		c.stratum.Collect(ch)
		c.rtt.Collect(ch)
		c.referenceTime.Collect(ch)
		c.rootDelay.Collect(ch)
		c.rootDistance.Collect(ch)
		c.rootDispersion.Collect(ch)
		c.precision.Collect(ch)
		c.leap.Collect(ch)
		c.scrapeDuration.Collect(ch)
	} else {
		log.Println("ERROR:", err)
		return
	}
}

func (c Collector) measure() error {
	const highDrift = 0.01

	begin := time.Now()
	measurement, err := c.getClockOffsetAndStratum()

	if err != nil {
		return fmt.Errorf("couldn't get NTP measurement: %s", err)
	}

	//if clock drift is unusually high (e.g. >10ms): repeat measurements for 30 seconds and submit median value
	if measurement.clockOffset > highDrift {
		//arrays of measurements used to calculate median
		var measurementsClockOffset []float64
		var measurementsStratum []float64
		var measurementsRTT []float64
		var measurementsReferenceTime []float64
		var measurementsRootDelay []float64
		var measurementsRootDispersion []float64
		var measurementsRootDistance []float64
		var measurementsPrecision []float64
		var measurementsLeap []float64

		log.Printf("WARN: clock drift is above %.2fs, taking multiple measurements for %.2f seconds", highDrift, c.NtpMeasurementDuration.Seconds())
		for time.Since(begin) < c.NtpMeasurementDuration {
			nextMeasurement, err := c.getClockOffsetAndStratum()
			if err != nil {
				return fmt.Errorf("couldn't get NTP measurement: %s", err)
			}

			measurementsClockOffset = append(measurementsClockOffset, nextMeasurement.clockOffset)
			measurementsStratum = append(measurementsStratum, nextMeasurement.stratum)
			measurementsRTT = append(measurementsRTT, nextMeasurement.rtt)
			measurementsReferenceTime = append(measurementsReferenceTime, nextMeasurement.referenceTime)
			measurementsRootDelay = append(measurementsRootDelay, nextMeasurement.rootDelay)
			measurementsRootDispersion = append(measurementsRootDispersion, nextMeasurement.rootDispersion)
			measurementsRootDistance = append(measurementsRootDistance, nextMeasurement.rootDistance)
			measurementsPrecision = append(measurementsPrecision, nextMeasurement.precision)
			measurementsLeap = append(measurementsLeap, nextMeasurement.leap)
		}

		measurement.clockOffset = calculateMedian(measurementsClockOffset)
		measurement.stratum = calculateMedian(measurementsStratum)
		measurement.rtt = calculateMedian(measurementsRTT)
		measurement.referenceTime = calculateMedian(measurementsReferenceTime)
		measurement.rootDelay = calculateMedian(measurementsRootDelay)
		measurement.rootDispersion = calculateMedian(measurementsRootDispersion)
		measurement.rootDistance = calculateMedian(measurementsRootDistance)
		measurement.precision = calculateMedian(measurementsPrecision)
		measurement.leap = calculateMedian(measurementsLeap)
	}

	c.drift.WithLabelValues(c.NtpServer).Set(measurement.clockOffset)
	c.stratum.WithLabelValues(c.NtpServer).Set(measurement.stratum)
	c.rtt.WithLabelValues(c.NtpServer).Set(measurement.rtt)
	c.referenceTime.WithLabelValues(c.NtpServer).Set(measurement.referenceTime)
	c.rootDelay.WithLabelValues(c.NtpServer).Set(measurement.rootDelay)
	c.rootDispersion.WithLabelValues(c.NtpServer).Set(measurement.rootDispersion)
	c.rootDistance.WithLabelValues(c.NtpServer).Set(measurement.rootDistance)
	c.precision.WithLabelValues(c.NtpServer).Set(measurement.precision)
	c.leap.WithLabelValues(c.NtpServer).Set(measurement.leap)

	c.scrapeDuration.Observe(time.Since(begin).Seconds())
	return nil
}

func (c Collector) getClockOffsetAndStratum() (measurement Measurement, err error) {
	options := ntp.QueryOptions{Version: c.NtpProtocolVersion}
	resp, err := ntp.QueryWithOptions(c.NtpServer, options)
	if err != nil {
		return measurement, fmt.Errorf("couldn't get NTP measurement: %s", err)
	}
	measurement.clockOffset = resp.ClockOffset.Seconds()
	measurement.stratum = float64(resp.Stratum)
	measurement.rtt = resp.RTT.Seconds()
	measurement.referenceTime = float64(resp.ReferenceTime.Unix())
	measurement.rootDelay = resp.RootDelay.Seconds()
	measurement.rootDispersion = resp.RootDispersion.Seconds()
	measurement.rootDistance = resp.RootDistance.Seconds()
	measurement.precision = resp.Precision.Seconds()
	measurement.leap = float64(resp.Leap)
	return measurement, nil
}

func calculateMedian(slice []float64) (median float64) {
	sort.Float64s(slice)

	middle := len(slice) / 2
	median = slice[middle]
	if len(slice)%2 == 0 {
		median = (median + slice[middle-1]) / 2
	}
	return median
}
