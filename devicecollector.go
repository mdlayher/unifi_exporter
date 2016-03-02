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
	Devices          *prometheus.Desc
	AdoptedDevices   *prometheus.Desc
	UnadoptedDevices *prometheus.Desc

	UptimeSecondsTotal *prometheus.Desc

	WirelessReceivedBytesTotal    *prometheus.Desc
	WirelessTransmittedBytesTotal *prometheus.Desc

	WirelessReceivedPacketsTotal    *prometheus.Desc
	WirelessTransmittedPacketsTotal *prometheus.Desc
	WirelessTransmittedDroppedTotal *prometheus.Desc

	WiredReceivedBytesTotal    *prometheus.Desc
	WiredTransmittedBytesTotal *prometheus.Desc

	WiredReceivedPacketsTotal    *prometheus.Desc
	WiredTransmittedPacketsTotal *prometheus.Desc

	Stations      *prometheus.Desc
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
		Devices: prometheus.NewDesc(
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

		UptimeSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "uptime_seconds_total"),
			"Device uptime in seconds",
			labelsDevice,
			nil,
		),

		WirelessReceivedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_received_bytes_total"),
			"Number of bytes received wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessTransmittedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_transmitted_bytes_total"),
			"Number of bytes transmitted wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessReceivedPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_received_packets_total"),
			"Number of packets received wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessTransmittedPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_transmitted_packets_total"),
			"Number of packets transmitted wirelessly by devices",
			labelsDevice,
			nil,
		),

		WirelessTransmittedDroppedTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wireless_transmitted_packets_dropped_total"),
			"Number of packets which are dropped on wireless transmission by devices",
			labelsDevice,
			nil,
		),

		WiredReceivedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_received_bytes_total"),
			"Number of bytes received using wired interface by devices",
			labelsDevice,
			nil,
		),

		WiredTransmittedBytesTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_transmitted_bytes_total"),
			"Number of bytes transmitted using wired interface by devices",
			labelsDevice,
			nil,
		),

		WiredReceivedPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_received_packets_total"),
			"Number of packets received using wired interface by devices",
			labelsDevice,
			nil,
		),

		WiredTransmittedPacketsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "wired_transmitted_packets_total"),
			"Number of packets transmitted using wired interface by devices",
			labelsDevice,
			nil,
		),

		Stations: prometheus.NewDesc(
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
func (c *DeviceCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	for _, s := range c.sites {
		devices, err := c.c.Devices(s.Name)
		if err != nil {
			return c.Devices, err
		}

		ch <- prometheus.MustNewConstMetric(
			c.Devices,
			prometheus.GaugeValue,
			float64(len(devices)),
			s.Description,
		)

		c.collectDeviceAdoptions(ch, s.Description, devices)
		c.collectDeviceUptime(ch, s.Description, devices)
		c.collectDeviceBytes(ch, s.Description, devices)
		c.collectDeviceStations(ch, s.Description, devices)
	}

	return nil, nil
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
			c.UptimeSecondsTotal,
			prometheus.CounterValue,
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
			c.WirelessReceivedBytesTotal,
			prometheus.CounterValue,
			float64(d.Stats.All.ReceiveBytes),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WirelessTransmittedBytesTotal,
			prometheus.CounterValue,
			float64(d.Stats.All.TransmitBytes),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.WirelessReceivedPacketsTotal,
			prometheus.CounterValue,
			float64(d.Stats.All.ReceivePackets),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WirelessTransmittedPacketsTotal,
			prometheus.CounterValue,
			float64(d.Stats.All.TransmitPackets),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WirelessTransmittedDroppedTotal,
			prometheus.CounterValue,
			float64(d.Stats.All.TransmitDropped),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.WiredReceivedBytesTotal,
			prometheus.CounterValue,
			float64(d.Stats.Uplink.ReceiveBytes),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WiredTransmittedBytesTotal,
			prometheus.CounterValue,
			float64(d.Stats.Uplink.TransmitBytes),
			labels...,
		)

		ch <- prometheus.MustNewConstMetric(
			c.WiredReceivedPacketsTotal,
			prometheus.CounterValue,
			float64(d.Stats.Uplink.ReceivePackets),
			labels...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.WiredTransmittedPacketsTotal,
			prometheus.CounterValue,
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
				c.Stations,
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
		c.Devices,
		c.AdoptedDevices,
		c.UnadoptedDevices,

		c.UptimeSecondsTotal,

		c.WirelessReceivedBytesTotal,
		c.WirelessTransmittedBytesTotal,

		c.WirelessReceivedPacketsTotal,
		c.WirelessTransmittedPacketsTotal,
		c.WirelessTransmittedDroppedTotal,

		c.WiredReceivedBytesTotal,
		c.WiredTransmittedBytesTotal,

		c.WiredReceivedPacketsTotal,
		c.WiredTransmittedPacketsTotal,

		c.Stations,
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
	if desc, err := c.collect(ch); err != nil {
		log.Printf("[ERROR] failed collecting device metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}
