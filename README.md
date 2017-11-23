unifi_exporter [![GoDoc](http://godoc.org/github.com/mdlayher/unifi_exporter?status.svg)](http://godoc.org/github.com/mdlayher/unifi_exporter) [![Build Status](https://travis-ci.org/mdlayher/unifi_exporter.svg?branch=master)](https://travis-ci.org/mdlayher/unifi_exporter) [![Coverage Status](https://coveralls.io/repos/mdlayher/unifi_exporter/badge.svg?branch=master)](https://coveralls.io/r/mdlayher/unifi_exporter?branch=master)
==============

Command `unifi_exporter` provides a Prometheus exporter for a Ubiquiti UniFi
Controller API and UniFi devices.

Package `unifiexporter` provides the Exporter type used in the `unifi_exporter`
Prometheus exporter.

MIT Licensed.

Usage
-----

```
$ ./unifi_exporter -h
Usage of ./unifi_exporter:
  -config.file string
       Relative path to config file yaml
```

To run the exporter, edit the included config.yml.example, rename it to config.yml, then run the exporter like so:

```
$ ./unifi_exporter -config.file config.yml
2017/11/15 17:06:32 [INFO] successfully authenticated to UniFi controller
2017/11/15 17:06:32 Starting UniFi exporter on ":9130" for site(s): Default
```

The minimum you'll need to modify is the unifi address, username and password. The port defaults to 8443 as specified in the config file,
and the defaults in 'listen' are sufficient for most users.

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
