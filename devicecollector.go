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
	TotalDevices     *prometheus.Desc
	AdoptedDevices   *prometheus.Desc
	UnadoptedDevices *prometheus.Desc

	UptimeSeconds *prometheus.Desc

	WirelessReceivedBytes    *prometheus.Desc
	WirelessTransmittedBytes *prometheus.Desc

	WirelessReceivedPackets    *prometheus.Desc
	WirelessTransmittedPackets *prometheus.Desc
	WirelessTransmittedDropped *prometheus.Desc

	WiredReceivedBytes    *prometheus.Desc
	WiredTransmittedBytes *prometheus.Desc

	WiredReceivedPackets    *prometheus.Desc
	WiredTransmittedPackets *prometheus.Desc

	TotalStations *prometheus.Desc
	UserStations  *prometheus.Desc
	GuestStations *prometheus.Desc

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
		TotalDevices: prometheus.NewDesc(
			// Subsystem is used as name so we get "unifi_devices"
			prometheus.BuildFQName(namespace, "", subsystem),
			"Total number of devices",
			labelsSiteOnly,
			nil,
		),

		AdoptedDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "adopted"),
			"Number of devices which are adopted",
			labelsSiteOnly,
			nil,
		),

		UnadoptedDevices: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "unadopted"),
			"Number of devices which are not adopted",
			labelsSiteOnly,
			nil,
		),

		UptimeSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "uptime_seconds"),
			"Device uptime in seconds",
			labelsDevice,
			nil,
		),

		WirelessReceivedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_received_bytes"),
			"Number of bytes received wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessTransmittedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_transmitted_bytes"),
			"Number of bytes transmitted wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessReceivedPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_received_packets"),
			"Number of packets received wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessTransmittedPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_transmitted_packets"),
			"Number of packets transmitted wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessTransmittedDropped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_transmitted_packets_dropped"),
			"Number of packets which are dropped on wireless transmission by devices",
			labelsDevice,
			nil,
		),

		WiredReceivedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_received_bytes"),
			"Number of bytes received using wired interface by devices",
			labelsDevice,
			nil,
		),

		WiredTransmittedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_transmitted_bytes"),
			"Number of bytes transmitted using wired interface by devices",
			labelsDevice,
			nil,
		),

		WiredReceivedPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_received_packets"),
			"Number of packets received using wired interface by devices",
			labelsDevice,
			nil,
		),

		WiredTransmittedPackets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_transmitted_packets"),
			"Number of packets transmitted using wired interface by devices",
			labelsDevice,
			nil,
		),

		TotalStations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "stations"),
			"Total number of stations (clients) connected to devices",
			labelsDeviceStations,
			nil,
		),

		UserStations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "stations_user"),
			"Number of user stations (private clients) connected to devices",
			labelsDeviceStations,
			nil,
		),

		GuestStations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "stations_guest"),
			"Number of guest stations (public clients) connected to devices",
			labelsDeviceStations,
			nil,
		),

		c:     c,
		sites: sites,
	}
}

// collect begins a metrics collection task for all metrics related to UniFi
// devices.
func (c *DeviceCollector) collect(ch chan<- prometheus.Metric) error {
	for _, s := range c.sites {
		devices, err := c.c.Devices(s.Name)
		if err != nil {
			return err
		}

		ch <- prometheus.MustNewConstMetric(
			c.TotalDevices,
			prometheus.GaugeValue,
			float64(len(devices)),
			s.Description,
		)

		c.collectDeviceAdoptions(ch, s.Description, devices)
		c.collectDeviceUptime(ch, s.Description, devices)
		c.collectDeviceBytes(ch, s.Description, devices)
		c.collectDeviceStations(ch, s.Description, devices)
	}

	return nil
}

// collectDeviceAdoptions collects counts for number of adopted and unadopted
// UniFi devices.
func (c *DeviceCollector) collectDeviceAdoptions(ch chan<- prometheus.Metric, siteLabel string, devices []*unifi.Device) {
	var adopted, unadopted int

	for _, d := range devices {
		if d.Adopted {
			adopted++
		} else {
			unadopted++
		}
	}

	ch <- prometheus.MustNewConstMetric(
		c.AdoptedDevices,
		prometheus.GaugeValue,
		float64(adopted),
		siteLabel,
	)

	ch <- prometheus.MustNewConstMetric(
		c.UnadoptedDevices,
		prometheus.GaugeValue,
		float64(unadopted),
		siteLabel,
	)
}

// collectDeviceUptime collects device uptime for UniFi devices.
func (c *DeviceCollector) collectDeviceUptime(ch chan<- prometheus.Metric, siteLabel string, devices []*unifi.Device) {
	for _, d := range devices {
		labels := []string{
			siteLabel,
			d.ID,
			d.NICs[0].MAC.String(),
			d.Name,
		}

		ch <- prometheus.MustNewConstMetric(
			c.UptimeSeconds,
			prometheus.GaugeValue,
			float64(d.Uptime/time.Second),
			labels...,
		)
	}
}

// collectDeviceBytes collects receive and transmit byte counts for UniFi devices.
func (c *DeviceCollector) collectDeviceBytes(ch chan<- prometheus.Metric, siteLabel string, devices []*unifi.Device) {
	for _, d := range devices {
		labels := []string{
			siteLabel,
			d.ID,
			d.NICs[0].MAC.String(),
			d.Name,
		}

		ch <- prometheus.MustNewConstMetric(
			c.WirelessReceivedBytes,
			prometheus.GaugeValue,
			float64(d.Stats.All.ReceiveBytes),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WirelessTransmittedBytes,
			prometheus.GaugeValue,
			float64(d.Stats.All.TransmitBytes),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.WirelessReceivedPackets,
			prometheus.GaugeValue,
			float64(d.Stats.All.ReceivePackets),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WirelessTransmittedPackets,
			prometheus.GaugeValue,
			float64(d.Stats.All.TransmitPackets),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WirelessTransmittedDropped,
			prometheus.GaugeValue,
			float64(d.Stats.All.TransmitDropped),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.WiredReceivedBytes,
			prometheus.GaugeValue,
			float64(d.Stats.Uplink.ReceiveBytes),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WiredTransmittedBytes,
			prometheus.GaugeValue,
			float64(d.Stats.Uplink.TransmitBytes),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.WiredReceivedPackets,
			prometheus.GaugeValue,
			float64(d.Stats.Uplink.ReceivePackets),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WiredTransmittedPackets,
			prometheus.GaugeValue,
			float64(d.Stats.Uplink.TransmitPackets),
			labels...,
		)
	}
}

// collectDeviceStations collects station counts for UniFi devices.
func (c *DeviceCollector) collectDeviceStations(ch chan<- prometheus.Metric, siteLabel string, devices []*unifi.Device) {
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

			ch <- prometheus.MustNewConstMetric(
				c.TotalStations,
				prometheus.GaugeValue,
				float64(r.Stats.NumberStations),
				llabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.UserStations,
				prometheus.GaugeValue,
				float64(r.Stats.NumberUserStations),
				llabels...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.GuestStations,
				prometheus.GaugeValue,
				float64(r.Stats.NumberGuestStations),
				llabels...,
			)
		}
	}
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *DeviceCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.TotalDevices,
		c.AdoptedDevices,
		c.UnadoptedDevices,

		c.UptimeSeconds,

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

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric pertaining to the global
// cluster usage over to the provided prometheus Metric channel.
func (c *DeviceCollector) Collect(ch chan<- prometheus.Metric) {
	if err := c.collect(ch); err != nil {
		log.Println("[ERROR] failed collecting device metrics:", err)
		return
	}
}
