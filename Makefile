.PHONY: all 

build:
	go build ./cmd/unifi_exporter

docker:
	docker build -t mdlayher/unifi_exporter .
