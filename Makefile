.PHONY: all 

build:
	go build ./cmd/unifi_exporter

docker:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ./cmd/unifi_exporter
	docker build -t mdlayher/unifi_exporter .
