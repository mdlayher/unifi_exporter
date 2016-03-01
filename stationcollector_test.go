package unifiexporter

import (
	"regexp"
	"strings"
	"testing"

	"github.com/mdlayher/unifi"
)

func TestStationCollector(t *testing.T) {
	var tests = []struct {
		desc    string
		input   string
		sites   []*unifi.Site
		matches []*regexp.Regexp
	}{
		{
			desc: "one station, one site",
			input: strings.TrimSpace(`
{
	"data": [
		{
			"_id": "abcdef",
			"ap_mac": "a0:a0:a0:a0:a0:a0",
			"mac": "de:ad:be:ef:de:ad",
			"hostname": "foo",
			"rx_bytes": 10,
			"rx_packets": 1,
			"tx_bytes": 20,
			"tx_packets": 2
		}
	]
}
`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_stations_total{site="Default"} 1`),

				regexp.MustCompile(`unifi_stations_received_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 10`),
				regexp.MustCompile(`unifi_stations_transmitted_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 20`),

				regexp.MustCompile(`unifi_stations_received_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 1`),
				regexp.MustCompile(`unifi_stations_transmitted_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 2`),
			},
			sites: []*unifi.Site{{
				Name:        "default",
				Description: "Default",
			}},
		},
		{
			desc: "two stations, one site",
			input: strings.TrimSpace(`
{
	"data": [
		{
			"_id": "abcdef",
			"ap_mac": "a0:a0:a0:a0:a0:a0",
			"mac": "de:ad:be:ef:de:ad",
			"hostname": "foo",
			"rx_bytes": 10,
			"rx_packets": 1,
			"tx_bytes": 20,
			"tx_packets": 2
		},
		{
			"_id": "123456",
			"ap_mac": "a0:a0:a0:a0:a0:a0",
			"mac": "ab:ad:1d:ea:ab:ad",
			"hostname": "bar",
			"rx_bytes": 100,
			"rx_packets": 10,
			"tx_bytes": 200,
			"tx_packets": 20
		}
	]
}
`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_stations_total{site="Default"} 2`),

				regexp.MustCompile(`unifi_stations_received_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 10`),
				regexp.MustCompile(`unifi_stations_transmitted_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 20`),

				regexp.MustCompile(`unifi_stations_received_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 1`),
				regexp.MustCompile(`unifi_stations_transmitted_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 2`),

				regexp.MustCompile(`unifi_stations_received_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="bar",id="123456",site="Default",station_mac="ab:ad:1d:ea:ab:ad"} 100`),
				regexp.MustCompile(`unifi_stations_transmitted_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="bar",id="123456",site="Default",station_mac="ab:ad:1d:ea:ab:ad"} 200`),

				regexp.MustCompile(`unifi_stations_received_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="bar",id="123456",site="Default",station_mac="ab:ad:1d:ea:ab:ad"} 10`),
				regexp.MustCompile(`unifi_stations_transmitted_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="bar",id="123456",site="Default",station_mac="ab:ad:1d:ea:ab:ad"} 20`),
			},
			sites: []*unifi.Site{{
				Name:        "default",
				Description: "Default",
			}},
		},
		{
			desc: "two stations, two sites (same station, but this is okay for tests)",
			input: strings.TrimSpace(`
{
	"data": [
		{
			"_id": "abcdef",
			"ap_mac": "a0:a0:a0:a0:a0:a0",
			"mac": "de:ad:be:ef:de:ad",
			"hostname": "foo",
			"rx_bytes": 10,
			"rx_packets": 1,
			"tx_bytes": 20,
			"tx_packets": 2
		}
	]
}
`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_stations_total{site="Default"} 1`),

				regexp.MustCompile(`unifi_stations_received_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 10`),
				regexp.MustCompile(`unifi_stations_transmitted_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 20`),

				regexp.MustCompile(`unifi_stations_received_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 1`),
				regexp.MustCompile(`unifi_stations_transmitted_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Default",station_mac="de:ad:be:ef:de:ad"} 2`),

				regexp.MustCompile(`unifi_stations_total{site="Some Site"} 1`),

				regexp.MustCompile(`unifi_stations_received_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Some Site",station_mac="de:ad:be:ef:de:ad"} 10`),
				regexp.MustCompile(`unifi_stations_transmitted_bytes{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Some Site",station_mac="de:ad:be:ef:de:ad"} 20`),

				regexp.MustCompile(`unifi_stations_received_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Some Site",station_mac="de:ad:be:ef:de:ad"} 1`),
				regexp.MustCompile(`unifi_stations_transmitted_packets{ap_mac="a0:a0:a0:a0:a0:a0",hostname="foo",id="abcdef",site="Some Site",station_mac="de:ad:be:ef:de:ad"} 2`),
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

		out := testStationCollector(t, []byte(tt.input), tt.sites)

		for j, m := range tt.matches {
			t.Logf("\t[%02d:%02d] match: %s", i, j, m.String())

			if !m.Match(out) {
				t.Fatal("\toutput failed to match regex")
			}
		}
	}
}

func testStationCollector(t *testing.T, input []byte, sites []*unifi.Site) []byte {
	c, done := testUniFiClient(t, input)
	defer done()

	collector := NewStationCollector(
		c,
		sites,
	)

	return testCollector(t, collector)
}
