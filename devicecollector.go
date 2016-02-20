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

	c     *unifi.Client
	sites []*unifi.Site
}

// Verify that the Exporter implements the prometheus.Collector interface.
var _ prometheus.Collector = &DeviceCollector{}

// NewDeviceCollector creates a new DeviceCollector which collects metrics for
// a specified site.
func NewDeviceCollector(c *unifi.Client, sites []*unifi.Site) *DeviceCollector {
	const subsystem = "devices"
	labels := []string{"site"}

	return &DeviceCollector{
		TotalDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "total",
				Help:      "Total number of devices registered to UniFi Controller, partitioned by site",
			},
			labels,
		),

		AdoptedDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "adopted",
				Help:      "Number of devices known to UniFi Controller which are adopted, partitioned by site",
			},
			labels,
		),

		UnadoptedDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "unadopted",
				Help:      "Number of devices known to UniFi Controller which are not adopted, partitioned by site",
			},
			labels,
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

		lv := siteDescription(s.Description)

		c.TotalDevices.WithLabelValues(lv).Set(float64(len(devices)))
		c.collectDeviceAdoptions(lv, devices)
	}

	return nil
}

// collectDeviceAdoptions collects counts for number of adopted and unadopted
// UniFi devices.
func (c *DeviceCollector) collectDeviceAdoptions(lv string, devices []*unifi.Device) {
	var adopted, unadopted int

	for _, d := range devices {
		if d.Adopted {
			adopted++
		} else {
			unadopted++
		}
	}

	c.AdoptedDevices.WithLabelValues(lv).Set(float64(adopted))
	c.UnadoptedDevices.WithLabelValues(lv).Set(float64(unadopted))
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
