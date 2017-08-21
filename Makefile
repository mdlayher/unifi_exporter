GO := CGO_ENABLED=0 go
PACKAGES = $(shell go list ./... | grep -v /vendor/)

.PHONY: all 
all: fmt vet lint test build

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: lint
lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

.PHONY: test
test:
	@for PKG in $(PACKAGES); do go test -v -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

build:
	$(GO) build ./cmd/unifi_exporter

docker:
	GOOS=linux GOARCH=amd64 $(GO) build ./cmd/unifi_exporter
	docker build -t mdlayher/unifi_exporter .
