// Package unifiexporter provides the Exporter type used in the unifi_exporter
// Prometheus exporter.
package unifiexporter

import (
	"log"
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
	collectors []collector
	sites      []*unifi.Site
	clientFn   ClientFunc
}

// Verify that the Exporter implements the prometheus.Collector interface.
var _ prometheus.Collector = &Exporter{}

// collector is essentially a modified prometheus.Collector which can return
// errors used to reconfigure the application.
type collector interface {
	prometheus.Collector
	CollectError(chan<- prometheus.Metric) error
}

// A ClientFunc is a function which can return an authenticated UniFi client.
// A ClientFunc is invoked by an Exporter whenever authentication against a UniFi
// controller fails, such as when a user's privileges are revoked or the
// authenticated session times out.
type ClientFunc func() (*unifi.Client, error)

// New creates a new Exporter which collects metrics from one or mote sites.
func New(sites []*unifi.Site, fn ClientFunc) (*Exporter, error) {
	e := &Exporter{
		clientFn: fn,
		sites:    sites,
	}

	if err := e.initClient(); err != nil {
		return nil, err
	}

	return e, nil
}

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, cc := range e.collectors {
		cc.Describe(ch)
	}
}

// Collect sends the collected metrics from each of the collectors to
// prometheus. Collect could be called several times concurrently
// and thus its run is protected by a single mutex.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, cc := range e.collectors {
		if err := cc.CollectError(ch); err == nil {
			continue
		}

		if err := e.initClient(); err != nil {
			log.Printf("[ERROR] could not initialize UniFi client: %v", err)
			return
		}
	}
}

// initClient sets up collectors for the Exporter, authenticating against
// the UniFi controller with a fresh session before doing so.
//
// initClient must be called with e's mutex locked.
func (e *Exporter) initClient() error {
	c, err := e.clientFn()
	if err != nil {
		return err
	}

	e.collectors = []collector{
		NewDeviceCollector(c, e.sites),
		NewStationCollector(c, e.sites),
	}

	log.Println("[INFO] successfully authenticated to UniFi controller")
	return nil
}
