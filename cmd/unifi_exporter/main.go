// Command unifi_exporter provides a Prometheus exporter for a Ubiquiti UniFi
// Controller API and UniFi devices.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mdlayher/unifi"
	"github.com/mdlayher/unifi_exporter"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// userAgent is ther user agent reported to the UniFi Controller API.
	userAgent = "github.com/mdlayher/unifi_exporter"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9130", "host:port for UniFi exporter")
	metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")

	unifiAddr = flag.String("unifi.addr", "", "address of UniFi Controller API")
	username  = flag.String("unifi.username", "", "username for authentication against UniFi Controller API")
	password  = flag.String("unifi.password", "", "password for authentication against UniFi Controller API")

	site     = flag.String("unifi.site", "", "[optional] description of site to collect metrics for using UniFi Controller API; if none specified, all sites will be scraped")
	insecure = flag.Bool("unifi.insecure", false, "[optional] do not verify TLS certificate for UniFi Controller API (warning: please use carefully)")
	timeout  = flag.Duration("unifi.timeout", 5*time.Second, "[optional] timeout for UniFi Controller API requests")
)

func main() {
	flag.Parse()

	if *unifiAddr == "" {
		log.Fatal("address of UniFi Controller API must be specified with '-unifi.addr' flag")
	}
	if *username == "" {
		log.Fatal("username to authenticate to UniFi Controller API must be specified with '-unifi.username' flag")
	}
	if *password == "" {
		log.Fatal("password to authenticate to UniFi Controller API must be specified with '-unifi.password' flag")
	}

	clientFn := newClient(
		*unifiAddr,
		*username,
		*password,
		*insecure,
		*timeout,
	)
	c, err := clientFn()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	sites, err := c.Sites()
	if err != nil {
		log.Fatalf("failed to retrieve list of sites: %v", err)
	}

	useSites, err := pickSites(*site, sites)
	if err != nil {
		log.Fatalf("failed to select a site: %v", err)
	}

	e, err := unifiexporter.New(useSites, clientFn)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	prometheus.MustRegister(e)

	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("Starting UniFi exporter on %q for site(s): %s", *telemetryAddr, sitesString(useSites))

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("cannot start UniFi exporter: %s", err)
	}
}

// pickSites attempts to find a site with a description matching the value
// specified in choose.  If choose is empty, all sites are returned.
func pickSites(choose string, sites []*unifi.Site) ([]*unifi.Site, error) {
	if choose == "" {
		return sites, nil
	}

	var pick *unifi.Site
	for _, s := range sites {
		if s.Description == choose {
			pick = s
			break
		}
	}
	if pick == nil {
		return nil, fmt.Errorf("site with description %q was not found in UniFi Controller", choose)
	}

	return []*unifi.Site{pick}, nil
}

// sitesString returns a comma-separated string of site descriptions, meant
// for displaying to users.
func sitesString(sites []*unifi.Site) string {
	ds := make([]string, 0, len(sites))
	for _, s := range sites {
		ds = append(ds, s.Description)
	}

	return strings.Join(ds, ", ")
}

// newClient returns a unifiexporter.ClientFunc using the input parameters.
func newClient(addr, username, password string, insecure bool, timeout time.Duration) unifiexporter.ClientFunc {
	return func() (*unifi.Client, error) {
		httpClient := &http.Client{Timeout: timeout}
		if insecure {
			httpClient = unifi.InsecureHTTPClient(timeout)
		}

		c, err := unifi.NewClient(addr, httpClient)
		if err != nil {
			return nil, fmt.Errorf("cannot create UniFi Controller client: %v", err)
		}
		c.UserAgent = userAgent

		if err := c.Login(username, password); err != nil {
			return nil, fmt.Errorf("failed to authenticate to UniFi Controller: %v", err)
		}

		return c, nil
	}
}
