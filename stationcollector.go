package unifiexporter

import (
	"log"

	"github.com/mdlayher/unifi"
	"github.com/prometheus/client_golang/prometheus"
)

// A StationCollector is a Prometheus collector for metrics regarding Ubiquiti
// UniFi stations (clients).
type StationCollector struct {
	Stations *prometheus.Desc

	ReceivedBytesTotal    *prometheus.Desc
	TransmittedBytesTotal *prometheus.Desc

	ReceivedPacketsTotal    *prometheus.Desc
	TransmittedPacketsTotal *prometheus.Desc

	RSSIDBM  *prometheus.Desc
	NoiseDBM *prometheus.Desc

	c     *unifi.Client
	sites []*unifi.Site
}

// Verify that the Exporter implements the prometheus.Collector interface.
var _ collector = &StationCollector{}

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
		Stations: prometheus.NewDesc(
			// Subsystem is used as name so we get "unifi_stations"
			prometheus.BuildFQName(namespace, "", subsystem),
			"Total number of stations (clients)",
			labelsSiteOnly,
			nil,
		),

		ReceivedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "received_bytes_total"),
			"Number of bytes received by stations (client download)",
			labelsStation,
			nil,
		),

		TransmittedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transmitted_bytes_total"),
			"Number of bytes transmitted by stations (client upload)",
			labelsStation,
			nil,
		),

		ReceivedPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "received_packets_total"),
			"Number of packets received by stations (client download)",
			labelsStation,
			nil,
		),

		TransmittedPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "transmitted_packets_total"),
			"Number of packets transmitted by stations (client upload)",
			labelsStation,
			nil,
		),

		RSSIDBM: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "rssi_dbm"),
			"Current signal strength of stations",
			labelsStation,
			nil,
		),

		NoiseDBM: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "noise_dbm"),
			"Current noise floor of stations",
			labelsStation,
			nil,
		),

		c:     c,
		sites: sites,
	}
}

// collect begins a metrics collection task for all metrics related to UniFi
// stations.
func (c *StationCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	for _, s := range c.sites {
		stations, err := c.c.Stations(s.Name)
		if err != nil {
			return c.Stations, err
		}

		ch <- prometheus.MustNewConstMetric(
			c.Stations,
			prometheus.GaugeValue,
			float64(len(stations)),
			s.Description,
		)

		c.collectStationBytes(ch, s.Description, stations)
		c.collectStationSignal(ch, s.Description, stations)
	}

	return nil, nil
}

// hostName picks the more desirable of the two names available. It uses the Unifi-set name if provided,
// otherwise uses the host-provided name.
func hostName(s *unifi.Station) string {
	if s.Name != "" {
		return s.Name
	} else {
		return s.Hostname
	}
}

// collectStationBytes collects receive and transmit byte counts for UniFi stations.
func (c *StationCollector) collectStationBytes(ch chan<- prometheus.Metric, siteLabel string, stations []*unifi.Station) {
	for _, s := range stations {
		labels := []string{
			siteLabel,
			s.ID,
			s.APMAC.String(),
			s.MAC.String(),
			hostName(s),
		}

		ch <- prometheus.MustNewConstMetric(
			c.ReceivedBytesTotal,
			prometheus.CounterValue,
			float64(s.Stats.ReceiveBytes),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TransmittedBytesTotal,
			prometheus.CounterValue,
			float64(s.Stats.TransmitBytes),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.ReceivedPacketsTotal,
			prometheus.CounterValue,
			float64(s.Stats.ReceivePackets),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.TransmittedPacketsTotal,
			prometheus.CounterValue,
			float64(s.Stats.TransmitPackets),
			labels...,
		)
	}
}

// collectStationSignal collects wireless signal strength for UniFi stations.
func (c *StationCollector) collectStationSignal(ch chan<- prometheus.Metric, siteLabel string, stations []*unifi.Station) {
	for _, s := range stations {
		if s.APMAC.String() == "" {
			continue
		}
		labels := []string{
			siteLabel,
			s.ID,
			s.APMAC.String(),
			s.MAC.String(),
			hostName(s),
		}

		ch <- prometheus.MustNewConstMetric(
			c.RSSIDBM,
			prometheus.GaugeValue,
			float64(s.RSSI),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.NoiseDBM,
			prometheus.GaugeValue,
			float64(s.Noise),
			labels...,
		)
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *StationCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.Stations,

		c.ReceivedBytesTotal,
		c.TransmittedBytesTotal,

		c.ReceivedPacketsTotal,
		c.TransmittedPacketsTotal,

		c.RSSIDBM,
		c.NoiseDBM,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect is the same as Collect, but ignores any errors which occur.
// Collect exists to satisfy the prometheus.Collector interface.
func (c *StationCollector) Collect(ch chan<- prometheus.Metric) {
	_ = c.CollectError(ch)
}

// CollectError sends the metric values for each metric pertaining to the global
// cluster usage over to the provided prometheus Metric channel, returning any
// errors which occur.
func (c *StationCollector) CollectError(ch chan<- prometheus.Metric) error {
	if desc, err := c.collect(ch); err != nil {
		log.Printf("[ERROR] failed collecting station metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return err
	}

	return nil
}
