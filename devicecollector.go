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

	WirelessTotalBytes       *prometheus.GaugeVec
	WirelessReceivedBytes    *prometheus.GaugeVec
	WirelessTransmittedBytes *prometheus.GaugeVec

	WirelessReceivedPackets    *prometheus.GaugeVec
	WirelessTransmittedPackets *prometheus.GaugeVec
	WirelessTransmittedDropped *prometheus.GaugeVec

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
				Help:      "Total number of devices, partitioned by site",
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

		WirelessTotalBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_total_bytes",
				Help:      "Total number of bytes received and transmitted wirelessly by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		WirelessReceivedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_received_bytes",
				Help:      "Number of bytes received wirelessly by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		WirelessTransmittedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_transmitted_bytes",
				Help:      "Number of bytes transmitted wirelessly by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		WirelessReceivedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_received_packets",
				Help:      "Number of packets received wirelessly by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		WirelessTransmittedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_transmitted_packets",
				Help:      "Number of packets transmitted wirelessly by devices, partitioned by site and device ID",
			},
			[]string{labelSite, labelID},
		),

		WirelessTransmittedDropped: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_transmitted_dropped",
				Help:      "Number of packets which are dropped on wireless transmission by devices, partitioned by site and device ID",
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

		c.WirelessTotalBytes,
		c.WirelessReceivedBytes,
		c.WirelessTransmittedBytes,

		c.WirelessReceivedPackets,
		c.WirelessTransmittedPackets,
		c.WirelessTransmittedDropped,
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
		c.WirelessTotalBytes.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.TotalBytes))
		c.WirelessReceivedBytes.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.ReceiveBytes))
		c.WirelessTransmittedBytes.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.TransmitBytes))

		c.WirelessReceivedPackets.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.ReceivePackets))
		c.WirelessTransmittedPackets.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.TransmitPackets))
		c.WirelessTransmittedDropped.WithLabelValues(siteLabel, d.ID).Set(float64(d.Stats.All.TransmitDropped))
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
