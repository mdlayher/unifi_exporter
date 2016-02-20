// Package unifiexporter provides the Exporter type used in the unifi_exporter
// Prometheus exporter.
package unifiexporter

import (
	"strings"
	"sync"

	"github.com/mdlayher/unifi"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// namespace is the top-level namespace for this UniFi exporter.
	namespace = "unifi"
)

// An Exporter is a Prometheus exporter for Ubiquiti UniFi Controller API
// metrics.  It wraps all UniFi metrics collectors and provides a single global
// exporter which can serve metrics. It also ensures that the collection
// is done in a thread-safe manner, the necessary requirement stated by
// Prometheus. It implements the prometheus.Collector interface in order to
// register with Prometheus.
type Exporter struct {
	mu         sync.Mutex
	collectors []prometheus.Collector
}

// Verify that the Exporter implements the prometheus.Collector interface.
var _ prometheus.Collector = &Exporter{}

// New creates a new Exporter which collects metrics from one or mote sites.
func New(c *unifi.Client, sites []*unifi.Site) *Exporter {
	return &Exporter{
		collectors: []prometheus.Collector{
			NewDeviceCollector(c, sites),
		},
	}
}

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (c *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, cc := range c.collectors {
		cc.Describe(ch)
	}
}

// Collect sends the collected metrics from each of the collectors to
// prometheus. Collect could be called several times concurrently
// and thus its run is protected by a single mutex.
func (c *Exporter) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, cc := range c.collectors {
		cc.Collect(ch)
	}
}

// siteDescription a metric label value for Prometheus metrics by
// normalizing the a site description field.
func siteDescription(desc string) string {
	desc = strings.ToLower(desc)
	return strings.Map(func(r rune) rune {
		switch r {
		// TODO(mdlayher): figure out valid set of characters for description
		case ' ', '-', '_', '.':
			return -1
		}

		return r
	}, desc)
}
