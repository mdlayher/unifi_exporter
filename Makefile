.PHONY: all 

VERSION := 0.4.0

build:
	go build -o unifi_exporter_${VERSION}_linux_amd64 ./cmd/unifi_exporter
	GOARCH=arm go build -o unifi_exporter_${VERSION}_linux_arm ./cmd/unifi_exporter


docker:
	docker build -t mdlayher/unifi_exporter .
