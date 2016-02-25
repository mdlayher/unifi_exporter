package unifiexporter

import (
	"log"

	"github.com/mdlayher/unifi"
	"github.com/prometheus/client_golang/prometheus"
)

// A StationCollector is a Prometheus collector for metrics regarding Ubiquiti
// UniFi stations (clients).
type StationCollector struct {
	TotalStations *prometheus.GaugeVec

	ReceivedBytes    *prometheus.GaugeVec
	TransmittedBytes *prometheus.GaugeVec

	ReceivedPackets    *prometheus.GaugeVec
	TransmittedPackets *prometheus.GaugeVec

	c     *unifi.Client
	sites []*unifi.Site
}

// Verify that the Exporter implements the prometheus.Collector interface.
var _ prometheus.Collector = &StationCollector{}

// NewStationCollector creates a new StationCollector which collects metrics for
// a specified site.
func NewStationCollector(c *unifi.Client, sites []*unifi.Site) *StationCollector {
	const (
		subsystem = "stations"
	)

	var (
		labelsSiteOnly = []string{"site"}
		labelsStation  = []string{"site", "id", "ap_mac", "station_mac", "hostname"}
	)

	return &StationCollector{
		TotalStations: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "total",
				Help:      "Total number of stations (clients), partitioned by site",
			},
			labelsSiteOnly,
		),

		ReceivedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "received_bytes",
				Help:      "Number of bytes received by stations (client download), partitioned by site, station, and access point",
			},
			labelsStation,
		),

		TransmittedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "transmitted_bytes",
				Help:      "Number of bytes transmitted by stations (client upload), partitioned by site, station, and access point",
			},
			labelsStation,
		),

		ReceivedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "received_packets",
				Help:      "Number of packets received by stations (client download), partitioned by site, station, and access point",
			},
			labelsStation,
		),

		TransmittedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "transmitted_packets",
				Help:      "Number of packets transmitted by stations (client upload), partitioned by site, station, and access point",
			},
			labelsStation,
		),

		c:     c,
		sites: sites,
	}
}

// collectors contains a list of collectors which are collected each time
// the exporter is scraped.  This list must be kept in sync with the collectors
// in StationCollector.
func (c *StationCollector) collectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.TotalStations,

		c.ReceivedBytes,
		c.TransmittedBytes,

		c.ReceivedPackets,
		c.TransmittedPackets,
	}
}

// collect begins a metrics collection task for all metrics related to UniFi
// stations.
func (c *StationCollector) collect() error {
	for _, s := range c.sites {
		stations, err := c.c.Stations(s.Name)
		if err != nil {
			return err
		}

		siteLabel := siteDescription(s.Description)

		c.TotalStations.WithLabelValues(siteLabel).Set(float64(len(stations)))
		c.collectStationBytes(siteLabel, stations)
	}

	return nil
}

// collectStationBytes collects receive and transmit byte counts for UniFi stations.
func (c *StationCollector) collectStationBytes(siteLabel string, stations []*unifi.Station) {
	for _, s := range stations {
		labels := []string{
			siteLabel,
			s.ID,
			s.APMAC.String(),
			s.MAC.String(),
			s.Hostname,
		}

		c.ReceivedBytes.WithLabelValues(labels...).Set(float64(s.Stats.ReceiveBytes))
		c.TransmittedBytes.WithLabelValues(labels...).Set(float64(s.Stats.TransmitBytes))

		c.ReceivedPackets.WithLabelValues(labels...).Set(float64(s.Stats.ReceivePackets))
		c.TransmittedPackets.WithLabelValues(labels...).Set(float64(s.Stats.TransmitPackets))
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *StationCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.collectors() {
		m.Describe(ch)
	}
}

// Collect sends the metric values for each metric pertaining to the global
// cluster usage over to the provided prometheus Metric channel.
func (c *StationCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(); err != nil {
		log.Fatalf("[ERROR] failed collecting station metrics: %v", err)
		return
	}

	for _, m := range c.collectors() {
		m.Collect(ch)
	}
}
