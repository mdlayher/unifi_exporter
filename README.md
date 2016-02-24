unifi_exporter [![GoDoc](http://godoc.org/github.com/mdlayher/unifi_exporter?status.svg)](http://godoc.org/github.com/mdlayher/unifi_exporter) [![Build Status](https://travis-ci.org/mdlayher/unifi_exporter.svg?branch=master)](https://travis-ci.org/mdlayher/unifi_exporter) [![Coverage Status](https://coveralls.io/repos/mdlayher/unifi_exporter/badge.svg?branch=master)](https://coveralls.io/r/mdlayher/unifi_exporter?branch=master)
==============

Command `unifi_exporter` provides a Prometheus exporter for a Ubiquiti UniFi
Controller API and UniFi devices.

Package `unifiexporter` provides the Exporter type used in the `unifi_exporter`
Prometheus exporter.

MIT Licensed.

Usage
-----

Available flags for `unifi_exporter` include:

```
$ ./unifi_exporter -h
Usage of ./unifi_exporter:
  -telemetry.addr string
        host:port for UniFi exporter (default ":9130")
  -telemetry.path string
        URL path for surfacing collected metrics (default "/metrics")
  -unifi.addr string
        address of UniFi Controller API
  -unifi.insecure
        [optional] do not verify TLS certificate for UniFi Controller API (warning: please use carefully)
  -unifi.password string
        password for authentication against UniFi Controller API
  -unifi.site string
        [optional] description of site to collect metrics for using UniFi Controller API; if none specified, all sites will be scraped
  -unifi.timeout duration
        [optional] timeout for UniFi Controller API requests (default 5s)
  -unifi.username string
        username for authentication against UniFi Controller API
```

An example of using `unifi_exporter` with authentication:

```
$ ./unifi_exporter -unifi.addr https://unifi.example.com:8443/ -unifi.username admin -unifi.password password
2016/02/24 13:41:34 Starting UniFi exporter on ":9130" for site(s): Foo, Bar, Baz
```

Sample
------

Here is a screenshot of a sample dashboard created using [`grafana`](https://github.com/grafana/grafana)
with metrics from exported from `unifi_exporter`.

![sample](https://cloud.githubusercontent.com/assets/1926905/13296555/163b39f2-dafc-11e5-84ef-8b8f03872c84.png)


Thanks
------

Special thanks to [Vaibhav Bhembre](https://github.com/neurodrone) for his work
on [`ceph_exporter`](https://github.com/digitalocean/ceph_exporter).
`ceph_exporter`  was used frequently as a reference Prometheus exporter while
implementing `unifi_exporter`.
