package unifiexporter

import (
	"log"

	"github.com/mdlayher/unifi"
	"github.com/prometheus/client_golang/prometheus"
)

// A DeviceCollector is a Prometheus collector for metrics regarding Ubiquiti
// UniFi devices.
type DeviceCollector struct {
	TotalDevices     *prometheus.GaugeVec
	AdoptedDevices   *prometheus.GaugeVec
	UnadoptedDevices *prometheus.GaugeVec

	TotalBytes       *prometheus.GaugeVec
	ReceivedBytes    *prometheus.GaugeVec
	TransmittedBytes *prometheus.GaugeVec

	ReceivedPackets    *prometheus.GaugeVec
	TransmittedPackets *prometheus.GaugeVec
	TransmittedDropped *prometheus.GaugeVec

	c     *unifi.Client
	sites []*unifi.Site
}

// Verify that the Exporter implements the prometheus.Collector interface.
var _ prometheus.Collector = &DeviceCollector{}

// NewDeviceCollector creates a new DeviceCollector which collects metrics for
// a specified site.
func NewDeviceCollector(c *unifi.Client, sites []*unifi.Site) *DeviceCollector {
	const (
		subsystem = "devices"

		labelSite = "site"
		labelID   = "id"
	)

	return &DeviceCollector{
		TotalDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "total",
				Help:      "Total number of devices registered , partitioned by site",
			},
			[]string{labelSite},
		),

		AdoptedDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "adopted",
				Help:      "Number of devices which are adopted, partitioned by site",
			},
			[]string{labelSite},
		),

		UnadoptedDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "unadopted",
				Help:      "Number of devices which are not adopted, partitioned by site",
			},
			[]string{labelSite},
		),

		TotalBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "total_bytes",
				Help:      "Total number of bytes received and transmitted by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		ReceivedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "received_bytes",
				Help:      "Number of bytes received by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		TransmittedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "transmitted_bytes",
				Help:      "Number of bytes transmitted by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		ReceivedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "received_packets",
				Help:      "Number of packets received by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		TransmittedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "transmitted_packets",
				Help:      "Number of packets transmitted by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		TransmittedDropped: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "transmitted_dropped",
				Help:      "Number of packets which are dropped on transmission by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		c:     c,
		sites: sites,
	}
}

// collectors contains a list of collectors which are collected each time
// the exporter is scraped.  This list must be kept in sync with the collectors
// in DeviceCollector.
func (c *DeviceCollector) collectors() []prometheus.Collector {
	return []prometheus.Collector{
		c.TotalDevices,
		c.AdoptedDevices,
		c.UnadoptedDevices,

		c.TotalBytes,
		c.ReceivedBytes,
		c.TransmittedBytes,

		c.ReceivedPackets,
		c.TransmittedPackets,
		c.TransmittedDropped,
	}
}

// collect begins a metrics collection task for all metrics related to UniFi
// devices.
func (c *DeviceCollector) collect() error {
	for _, s := range c.sites {
		devices, err := c.c.Devices(s.Name)
		if err != nil {
			return err
		}

		siteLabel := siteDescription(s.Description)

		c.TotalDevices.WithLabelValues(siteLabel).Set(float64(len(devices)))
		c.collectDeviceAdoptions(siteLabel, devices)
		c.collectDeviceBytes(siteLabel, devices)
	}

	return nil
}

// collectDeviceAdoptions collects counts for number of adopted and unadopted
// UniFi devices.
func (c *DeviceCollector) collectDeviceAdoptions(siteLabel string, devices []*unifi.Device) {
	var adopted, unadopted int

	for _, d := range devices {
		if d.Adopted {
			adopted++
		} else {
			unadopted++
		}
	}

	c.AdoptedDevices.WithLabelValues(siteLabel).Set(float64(adopted))
	c.UnadoptedDevices.WithLabelValues(siteLabel).Set(float64(unadopted))
}

// collectDeviceBytes collects receive and transmit byte counts for UniFi devices.
func (c *DeviceCollector) collectDeviceBytes(siteLabel string, devices []*unifi.Device) {
	for _, d := range devices {
		c.TotalBytes.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.TotalBytes))
		c.ReceivedBytes.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.ReceiveBytes))
		c.TransmittedBytes.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.TransmitBytes))

		c.ReceivedPackets.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.ReceivePackets))
		c.TransmittedPackets.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.TransmitPackets))
		c.TransmittedDropped.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.TransmitDropped))
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *DeviceCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.collectors() {
		m.Describe(ch)
	}
}

// Collect sends the metric values for each metric pertaining to the global
// cluster usage over to the provided prometheus Metric channel.
func (c *DeviceCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(); err != nil {
		log.Fatalf("[ERROR] failed collecting device metrics: %v", err)
		return
	}

	for _, m := range c.collectors() {
		m.Collect(ch)
	}
}
