package unifiexporter

import (
	"regexp"
	"strings"
	"testing"

	"github.com/mdlayher/unifi"
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
			"name": "ABC",
			"ethernet_table": [{
				"mac": "de:ad:be:ef:de:ad"
			}],
			"radio_table_stats": [{
					"guest-num_sta": 1,
					"name": "wifi0",
					"num_sta": 3,
					"user-num_sta": 2
				}, {
					"guest-num_sta": 2,
					"name": "wifi1",
					"num_sta": 6,
					"user-num_sta": 4
			}],
			"radio_table": [
				{
					"name": "wifi0",
					"radio": "ng"
				},
				{
					"name": "wifi1",
					"radio": "na"
				}
			],
			"stat": {
				"bytes": 100,
				"rx_bytes": 80,
				"tx_bytes": 20,
				"rx_packets": 4,
				"tx_packets": 1,
				"tx_dropped": 1
			},
			"uplink": {
				"rx_bytes": 20,
				"tx_bytes": 10,
				"rx_packets": 2,
				"tx_packets": 1
			},
			"uptime": 10
		}
	]
}
`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_devices{site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_adopted{site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="Default"} 0`),

				regexp.MustCompile(`unifi_devices_uptime_seconds_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 10`),

				regexp.MustCompile(`unifi_devices_wireless_received_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_dropped_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_wired_received_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 20`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 10`),

				regexp.MustCompile(`unifi_devices_wired_received_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_stations{id="abc",interface="wifi0",mac="de:ad:be:ef:de:ad",name="ABC",radio="2.4GHz",site="Default"} 3`),
				regexp.MustCompile(`unifi_devices_stations{id="abc",interface="wifi1",mac="de:ad:be:ef:de:ad",name="ABC",radio="5GHz",site="Default"} 6`),
				regexp.MustCompile(`unifi_devices_stations_user{id="abc",interface="wifi0",mac="de:ad:be:ef:de:ad",name="ABC",radio="2.4GHz",site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_stations_user{id="abc",interface="wifi1",mac="de:ad:be:ef:de:ad",name="ABC",radio="5GHz",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="abc",interface="wifi0",mac="de:ad:be:ef:de:ad",name="ABC",radio="2.4GHz",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="abc",interface="wifi1",mac="de:ad:be:ef:de:ad",name="ABC",radio="5GHz",site="Default"} 2`),
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
					"name": "ABC",
					"ethernet_table": [{
						"mac": "de:ad:be:ef:de:ad"
					}],
					"radio_table_stats": [{
						"guest-num_sta": 1,
						"name": "wifi0",
						"num_sta": 3,
						"user-num_sta": 2
					}, {
						"guest-num_sta": 2,
						"name": "wifi1",
						"num_sta": 6,
						"user-num_sta": 4
					}],
					"radio_table": [
						{
							"name": "wifi0",
							"radio": "ng"
						},
						{
							"name": "wifi1",
							"radio": "na"
						}
					],
					"stat": {
						"bytes": 100,
						"rx_bytes": 80,
						"tx_bytes": 20,
						"rx_packets": 4,
						"tx_packets": 1,
						"tx_dropped": 1
					},
					"uplink": {
						"rx_bytes": 20,
						"tx_bytes": 10,
						"rx_packets": 2,
						"tx_packets": 1
					},
					"uptime": 10
				},
				{
					"_id": "def",
					"adopted": false,
					"inform_ip": "192.168.1.1",
					"name": "DEF",
					"ethernet_table": [{
						"mac": "ab:ad:1d:ea:ab:ad"
					}],
					"radio_table_stats": [{
						"guest-num_sta": 1,
						"name": "wifi0",
						"num_sta": 3,
						"user-num_sta": 2
					}, {
						"guest-num_sta": 2,
						"name": "wifi1",
						"num_sta": 6,
						"user-num_sta": 4
					}],
					"radio_table": [
						{
							"name": "wifi0",
							"radio": "ng"
						},
						{
							"name": "wifi1",
							"radio": "na"
						}
					],
					"stat": {
						"bytes": 200,
						"rx_bytes": 10,
						"tx_bytes": 190,
						"rx_packets": 1,
						"tx_packets": 19,
						"tx_dropped": 1
					},
					"uplink": {
						"rx_bytes": 40,
						"tx_bytes": 20,
						"rx_packets": 4,
						"tx_packets": 2
					},
					"uptime": 20
				}
			]
		}
		`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_devices{site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_adopted{site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_uptime_seconds_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 10`),

				regexp.MustCompile(`unifi_devices_wireless_received_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_dropped_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_wired_received_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 20`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_bytes_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 10`),

				regexp.MustCompile(`unifi_devices_wired_received_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_packets_total{id="abc",mac="de:ad:be:ef:de:ad",name="ABC",site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_stations{id="abc",interface="wifi0",mac="de:ad:be:ef:de:ad",name="ABC",radio="2.4GHz",site="Default"} 3`),
				regexp.MustCompile(`unifi_devices_stations{id="abc",interface="wifi1",mac="de:ad:be:ef:de:ad",name="ABC",radio="5GHz",site="Default"} 6`),
				regexp.MustCompile(`unifi_devices_stations_user{id="abc",interface="wifi0",mac="de:ad:be:ef:de:ad",name="ABC",radio="2.4GHz",site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_stations_user{id="abc",interface="wifi1",mac="de:ad:be:ef:de:ad",name="ABC",radio="5GHz",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="abc",interface="wifi0",mac="de:ad:be:ef:de:ad",name="ABC",radio="2.4GHz",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="abc",interface="wifi1",mac="de:ad:be:ef:de:ad",name="ABC",radio="5GHz",site="Default"} 2`),

				regexp.MustCompile(`unifi_devices_uptime_seconds_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_bytes_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 10`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 190`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 19`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_dropped_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_wired_received_bytes_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 40`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_bytes_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 20`),

				regexp.MustCompile(`unifi_devices_wired_received_packets_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_packets_total{id="def",mac="ab:ad:1d:ea:ab:ad",name="DEF",site="Default"} 2`),

				regexp.MustCompile(`unifi_devices_stations{id="def",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="DEF",radio="2.4GHz",site="Default"} 3`),
				regexp.MustCompile(`unifi_devices_stations{id="def",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="DEF",radio="5GHz",site="Default"} 6`),
				regexp.MustCompile(`unifi_devices_stations_user{id="def",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="DEF",radio="2.4GHz",site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_stations_user{id="def",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="DEF",radio="5GHz",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="def",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="DEF",radio="2.4GHz",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="def",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="DEF",radio="5GHz",site="Default"} 2`),
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
					"name": "OneTwoThree",
					"ethernet_table": [{
						"mac": "ab:ad:1d:ea:ab:ad"
					}],
					"radio_table_stats": [{
						"guest-num_sta": 1,
						"name": "wifi0",
						"num_sta": 3,
						"user-num_sta": 2
					}, {
						"guest-num_sta": 2,
						"name": "wifi1",
						"num_sta": 6,
						"user-num_sta": 4
					}],
					"radio_table": [
						{
							"name": "wifi0",
							"radio": "ng"
						},
						{
							"name": "wifi1",
							"radio": "na"
						}
					],
					"stat": {
						"bytes": 100,
						"rx_bytes": 80,
						"tx_bytes": 20,
						"rx_packets": 4,
						"tx_packets": 1,
						"tx_dropped": 1
					},
					"uplink": {
						"rx_bytes": 20,
						"tx_bytes": 10,
						"rx_packets": 2,
						"tx_packets": 1
					},
					"uptime": 10
				}
			]
		}
		`),
			matches: []*regexp.Regexp{
				regexp.MustCompile(`unifi_devices{site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_adopted{site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="Default"} 0`),

				regexp.MustCompile(`unifi_devices_uptime_seconds_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 10`),

				regexp.MustCompile(`unifi_devices_wireless_received_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_dropped_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_wired_received_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 20`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 10`),

				regexp.MustCompile(`unifi_devices_wired_received_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Default"} 1`),

				regexp.MustCompile(`unifi_devices_stations{id="123",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="2.4GHz",site="Default"} 3`),
				regexp.MustCompile(`unifi_devices_stations{id="123",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="5GHz",site="Default"} 6`),
				regexp.MustCompile(`unifi_devices_stations_user{id="123",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="2.4GHz",site="Default"} 2`),
				regexp.MustCompile(`unifi_devices_stations_user{id="123",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="5GHz",site="Default"} 4`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="123",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="2.4GHz",site="Default"} 1`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="123",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="5GHz",site="Default"} 2`),

				regexp.MustCompile(`unifi_devices{site="Some Site"} 1`),
				regexp.MustCompile(`unifi_devices_adopted{site="Some Site"} 1`),
				regexp.MustCompile(`unifi_devices_unadopted{site="Some Site"} 0`),

				regexp.MustCompile(`unifi_devices_uptime_seconds_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 10`),

				regexp.MustCompile(`unifi_devices_wireless_received_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 80`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 20`),

				regexp.MustCompile(`unifi_devices_wireless_received_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 4`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 1`),
				regexp.MustCompile(`unifi_devices_wireless_transmitted_packets_dropped_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 1`),

				regexp.MustCompile(`unifi_devices_wired_received_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 20`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_bytes_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 10`),

				regexp.MustCompile(`unifi_devices_wired_received_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 2`),
				regexp.MustCompile(`unifi_devices_wired_transmitted_packets_total{id="123",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",site="Some Site"} 1`),

				regexp.MustCompile(`unifi_devices_stations{id="123",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="2.4GHz",site="Some Site"} 3`),
				regexp.MustCompile(`unifi_devices_stations{id="123",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="5GHz",site="Some Site"} 6`),
				regexp.MustCompile(`unifi_devices_stations_user{id="123",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="2.4GHz",site="Some Site"} 2`),
				regexp.MustCompile(`unifi_devices_stations_user{id="123",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="5GHz",site="Some Site"} 4`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="123",interface="wifi0",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="2.4GHz",site="Some Site"} 1`),
				regexp.MustCompile(`unifi_devices_stations_guest{id="123",interface="wifi1",mac="ab:ad:1d:ea:ab:ad",name="OneTwoThree",radio="5GHz",site="Some Site"} 2`),
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

		out := testDeviceCollector(t, []byte(tt.input), tt.sites)

		for j, m := range tt.matches {
			t.Logf("\t[%02d:%02d] match: %s", i, j, m.String())

			if !m.Match(out) {
				t.Fatal("\toutput failed to match regex.")
			}
		}
	}
}

func testDeviceCollector(t *testing.T, input []byte, sites []*unifi.Site) []byte {
	c, done := testUniFiClient(t, input)
	defer done()

	collector := NewDeviceCollector(
		c,
		sites,
	)

	return testCollector(t, collector)
}
