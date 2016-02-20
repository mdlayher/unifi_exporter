unifi_exporter [![GoDoc](http://godoc.org/github.com/mdlayher/unifi_exporter?status.svg)](http://godoc.org/github.com/mdlayher/unifi_exporter) [![Build Status](https://travis-ci.org/mdlayher/unifi_exporter.svg?branch=master)](https://travis-ci.org/mdlayher/unifi_exporter) [![Coverage Status](https://coveralls.io/repos/mdlayher/unifi_exporter/badge.svg?branch=master)](https://coveralls.io/r/mdlayher/unifi_exporter?branch=master) [![Report Card](http://goreportcard.com/badge/mdlayher/unifi_exporter)](http://goreportcard.com/report/mdlayher/unifi_exporter)
==============

Command `unifi_exporter` provides a Prometheus exporter for a Ubiquiti UniFi
Controller API and UniFi devices.

Package `unifiexporter` provides the Exporter type used in the `unifi_exporter`
Prometheus exporter.

MIT Licensed.

Thanks
------

Special thanks to [Vaibhav Bhembre](https://github.com/neurodrone) for his work
on [`ceph_exporter`](https://github.com/digitalocean/ceph_exporter).
`ceph_exporter`  was used frequently as a reference Prometheus exporter while
implementing `unifi_exporter`.
