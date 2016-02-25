package unifiexporter

import (
	"log"
	"time"

	"github.com/mdlayher/unifi"
	"github.com/prometheus/client_golang/prometheus"
)

// A DeviceCollector is a Prometheus collector for metrics regarding Ubiquiti
// UniFi devices.
type DeviceCollector struct {
	TotalDevices     *prometheus.GaugeVec
	AdoptedDevices   *prometheus.GaugeVec
	UnadoptedDevices *prometheus.GaugeVec

	UptimeSeconds *prometheus.GaugeVec

	WirelessTotalBytes       *prometheus.GaugeVec
	WirelessReceivedBytes    *prometheus.GaugeVec
	WirelessTransmittedBytes *prometheus.GaugeVec

	WirelessReceivedPackets    *prometheus.GaugeVec
	WirelessTransmittedPackets *prometheus.GaugeVec
	WirelessTransmittedDropped *prometheus.GaugeVec

	WiredReceivedBytes    *prometheus.GaugeVec
	WiredTransmittedBytes *prometheus.GaugeVec

	WiredReceivedPackets    *prometheus.GaugeVec
	WiredTransmittedPackets *prometheus.GaugeVec

	TotalStations *prometheus.GaugeVec
	UserStations  *prometheus.GaugeVec
	GuestStations *prometheus.GaugeVec

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
	)

	var (
		labelsSiteOnly       = []string{"site"}
		labelsDevice         = []string{"site", "id", "mac", "name"}
		labelsDeviceStations = []string{"site", "id", "mac", "name", "interface", "radio"}
	)

	return &DeviceCollector{
		TotalDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "total",
				Help:      "Total number of devices, partitioned by site",
			},
			labelsSiteOnly,
		),

		AdoptedDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "adopted",
				Help:      "Number of devices which are adopted, partitioned by site",
			},
			labelsSiteOnly,
		),

		UnadoptedDevices: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "unadopted",
				Help:      "Number of devices which are not adopted, partitioned by site",
			},
			labelsSiteOnly,
		),

		UptimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "uptime_seconds",
				Help:      "Device uptime in seconds, partitioned by site and device",
			},
			labelsDevice,
		),

		WirelessTotalBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_total_bytes",
				Help:      "Total number of bytes received and transmitted wirelessly by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WirelessReceivedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_received_bytes",
				Help:      "Number of bytes received wirelessly by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WirelessTransmittedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_transmitted_bytes",
				Help:      "Number of bytes transmitted wirelessly by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WirelessReceivedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_received_packets",
				Help:      "Number of packets received wirelessly by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WirelessTransmittedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_transmitted_packets",
				Help:      "Number of packets transmitted wirelessly by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WirelessTransmittedDropped: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wireless_transmitted_dropped",
				Help:      "Number of packets which are dropped on wireless transmission by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WiredReceivedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wired_received_bytes",
				Help:      "Number of bytes received using wired interface by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WiredTransmittedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wired_transmitted_bytes",
				Help:      "Number of bytes transmitted using wired interface by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WiredReceivedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wired_received_packets",
				Help:      "Number of packets received using wired interface by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		WiredTransmittedPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "wired_transmitted_packets",
				Help:      "Number of packets transmitted using wired interface by devices, partitioned by site and device",
			},
			labelsDevice,
		),

		TotalStations: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "stations_total",
				Help:      "Total number of stations (clients) connected to devices, partitioned by site, device, and wireless radio",
			},
			labelsDeviceStations,
		),

		UserStations: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "stations_user",
				Help:      "Number of user stations (private clients) connected to devices, partitioned by site, device, and wireless radio",
			},
			labelsDeviceStations,
		),

		GuestStations: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "stations_guest",
				Help:      "Number of guest stations (public clients) connected to devices, partitioned by site, device, and wireless radio",
			},
			labelsDeviceStations,
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

		c.UptimeSeconds,

		c.WirelessTotalBytes,
		c.WirelessReceivedBytes,
		c.WirelessTransmittedBytes,

		c.WirelessReceivedPackets,
		c.WirelessTransmittedPackets,
		c.WirelessTransmittedDropped,

		c.WiredReceivedBytes,
		c.WiredTransmittedBytes,

		c.WiredReceivedPackets,
		c.WiredTransmittedPackets,

		c.TotalStations,
		c.UserStations,
		c.GuestStations,
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
		c.collectDeviceUptime(siteLabel, devices)
		c.collectDeviceBytes(siteLabel, devices)
		c.collectDeviceStations(siteLabel, devices)
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

// collectDeviceUptime collects device uptime for UniFi devices.
func (c *DeviceCollector) collectDeviceUptime(siteLabel string, devices []*unifi.Device) {
	for _, d := range devices {
		labels := []string{
			siteLabel,
			d.ID,
			d.NICs[0].MAC.String(),
			d.Name,
		}

		c.UptimeSeconds.WithLabelValues(labels...).Set(float64(d.Uptime / time.Second))
	}
}

// collectDeviceBytes collects receive and transmit byte counts for UniFi devices.
func (c *DeviceCollector) collectDeviceBytes(siteLabel string, devices []*unifi.Device) {
	for _, d := range devices {
		labels := []string{
			siteLabel,
			d.ID,
			d.NICs[0].MAC.String(),
			d.Name,
		}

		c.WirelessTotalBytes.WithLabelValues(labels...).Set(float64(d.Stats.TotalBytes))
		c.WirelessReceivedBytes.WithLabelValues(labels...).Set(float64(d.Stats.All.ReceiveBytes))
		c.WirelessTransmittedBytes.WithLabelValues(labels...).Set(float64(d.Stats.All.TransmitBytes))

		c.WirelessReceivedPackets.WithLabelValues(labels...).Set(float64(d.Stats.All.ReceivePackets))
		c.WirelessTransmittedPackets.WithLabelValues(labels...).Set(float64(d.Stats.All.TransmitPackets))
		c.WirelessTransmittedDropped.WithLabelValues(labels...).Set(float64(d.Stats.All.TransmitDropped))

		c.WiredReceivedBytes.WithLabelValues(labels...).Set(float64(d.Stats.Uplink.ReceiveBytes))
		c.WiredTransmittedBytes.WithLabelValues(labels...).Set(float64(d.Stats.Uplink.TransmitBytes))

		c.WiredReceivedPackets.WithLabelValues(labels...).Set(float64(d.Stats.Uplink.ReceivePackets))
		c.WiredTransmittedPackets.WithLabelValues(labels...).Set(float64(d.Stats.Uplink.TransmitPackets))
	}
}

// collectDeviceStations collects station counts for UniFi devices.
func (c *DeviceCollector) collectDeviceStations(siteLabel string, devices []*unifi.Device) {
	for _, d := range devices {
		labels := []string{
			siteLabel,
			d.ID,
			d.NICs[0].MAC.String(),
			d.Name,
		}

		for _, r := range d.Radios {
			// Since the radio name and type will be different for each
			// radio, we copy the original labels slice and append, to avoid
			// mutating it
			llabels := make([]string, len(labels))
			copy(llabels, labels)
			llabels = append(llabels, r.Name, r.Radio)

			c.TotalStations.WithLabelValues(llabels...).Set(float64(r.Stats.NumberStations))
			c.UserStations.WithLabelValues(llabels...).Set(float64(r.Stats.NumberUserStations))
			c.GuestStations.WithLabelValues(llabels...).Set(float64(r.Stats.NumberGuestStations))
		}
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
