package unifiexporter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/mdlayher/unifi"
	"github.com/prometheus/client_golang/prometheus"
)

func TestDeviceCollector(t *testing.T) {
	var tests = []struct {
		desc    string
		input   string
		sites   []*unifi.Site
		matches []*regexp.Regexp
	}{
		{
			desc: "one device, one site",
			input: strings.TrimSpace(`
{
	"data": [
		{
			"_id": "abc",
			"adopted": true,
			"inform_ip": "192.168.1.1",
			"stat": {
				"bytes": 100,
				"rx_bytes": 80,
				"tx_bytes": 20,
				"rx_packets": 4,
				"tx_packets": 1,
				"tx_dropped": 1
			}
		}
	]
}
`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_devices_total{site="default"} 1`),
				regexp.MustCompile(`unifi_devices_adopted{site="default"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="default"} 0`),

				regexp.MustCompile(`unifi_devices_wireless_total_bytes{id="abc",site="default"} 100`),
				regexp.MustCompile(`unifi_devices_wireless_received_bytes{id="abc",site="default"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes{id="abc",site="default"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets{id="abc",site="default"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets{id="abc",site="default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_dropped{id="abc",site="default"} 1`),
			},
			sites: []*unifi.Site{{
				Name:        "default",
				Description: "Default",
			}},
		},
		{
			desc: "two devices, one site",
			input: strings.TrimSpace(`
{
	"data": [
		{
			"_id": "abc",
			"adopted": true,
			"inform_ip": "192.168.1.1",
			"stat": {
				"bytes": 100,
				"rx_bytes": 80,
				"tx_bytes": 20,
				"rx_packets": 4,
				"tx_packets": 1,
				"tx_dropped": 1
			}
		},
		{
			"_id": "def",
			"adopted": false,
			"inform_ip": "192.168.1.1",
			"stat": {
				"bytes": 200,
				"rx_bytes": 10,
				"tx_bytes": 190,
				"rx_packets": 1,
				"tx_packets": 19,
				"tx_dropped": 1
			}
		}
	]
}
`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_devices_total{site="default"} 2`),
				regexp.MustCompile(`unifi_devices_adopted{site="default"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="default"} 1`),

				regexp.MustCompile(`unifi_devices_wireless_total_bytes{id="abc",site="default"} 100`),
				regexp.MustCompile(`unifi_devices_wireless_received_bytes{id="abc",site="default"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes{id="abc",site="default"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets{id="abc",site="default"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets{id="abc",site="default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_dropped{id="abc",site="default"} 1`),

				regexp.MustCompile(`unifi_devices_wireless_total_bytes{id="def",site="default"} 200`),
				regexp.MustCompile(`unifi_devices_wireless_received_bytes{id="def",site="default"} 10`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes{id="def",site="default"} 190`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets{id="def",site="default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets{id="def",site="default"} 19`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_dropped{id="def",site="default"} 1`),
			},
			sites: []*unifi.Site{{
				Name:        "default",
				Description: "Default",
			}},
		},
		{
			desc: "two devices, two sites (same device, but this is okay for tests)",
			input: strings.TrimSpace(`
{
	"data": [
		{
			"_id": "123",
			"adopted": true,
			"inform_ip": "192.168.1.1",
			"stat": {
				"bytes": 100,
				"rx_bytes": 80,
				"tx_bytes": 20,
				"rx_packets": 4,
				"tx_packets": 1,
				"tx_dropped": 1
			}
		}
	]
}
`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_devices_total{site="default"} 1`),
				regexp.MustCompile(`unifi_devices_adopted{site="default"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="default"} 0`),

				regexp.MustCompile(`unifi_devices_wireless_total_bytes{id="123",site="default"} 100`),
				regexp.MustCompile(`unifi_devices_wireless_received_bytes{id="123",site="default"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes{id="123",site="default"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets{id="123",site="default"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets{id="123",site="default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_dropped{id="123",site="default"} 1`),

				regexp.MustCompile(`unifi_devices_total{site="somesite"} 1`),
				regexp.MustCompile(`unifi_devices_adopted{site="somesite"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="somesite"} 0`),

				regexp.MustCompile(`unifi_devices_wireless_total_bytes{id="123",site="somesite"} 100`),
				regexp.MustCompile(`unifi_devices_wireless_received_bytes{id="123",site="somesite"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes{id="123",site="somesite"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets{id="123",site="somesite"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets{id="123",site="somesite"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_dropped{id="123",site="somesite"} 1`),
			},
			sites: []*unifi.Site{
				{
					Name:        "default",
					Description: "Default",
				},
				{
					Name:        "abcdef",
					Description: "Some Site",
				},
			},
		},
	}

	for i, tt := range tests {
		t.Logf("[%02d] test %q", i, tt.desc)

		out := testCollector(t, []byte(tt.input), tt.sites)

		for j, m := range tt.matches {
			t.Logf("\t[%02d:%02d] match: %s", i, j, m.String())

			if !m.Match(out) {
				t.Fatal("\toutput failed to match regex")
			}
		}
	}
}

func testCollector(t *testing.T, input []byte, sites []*unifi.Site) []byte {
	unifiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		_, _ = w.Write(input)
	}))
	defer unifiServer.Close()

	c, err := unifi.NewClient(unifiServer.URL, nil)
	if err != nil {
		t.Fatalf("failed to create UniFi client: %v", err)
	}

	collector := NewDeviceCollector(c, sites)

	if err := prometheus.Register(collector); err != nil {
		t.Fatalf("failed to register Prometheus collector: %v", err)
	}
	defer prometheus.Unregister(collector)

	promServer := httptest.NewServer(prometheus.Handler())
	defer promServer.Close()

	resp, err := http.Get(promServer.URL)
	if err != nil {
		t.Fatalf("failed to GET data from prometheus: %v", err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read server response: %v", err)
	}

	return buf
}
