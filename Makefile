PREFIX = /usr

all: build/ntp_exporter

GO_BUILDFLAGS = -mod vendor
GO_LDFLAGS    = -s -w -X main.version=$(shell ./util/find_version.sh)

build/ntp_exporter: FORCE
	go build $(GO_BUILDFLAGS) -ldflags '-s -w $(GO_LDFLAGS)' -o $@ .

install: build/ntp_exporter FORCE
	install -D -m 0755 build/ntp_exporter "$(DESTDIR)$(PREFIX)/bin/ntp_exporter"

vendor: FORCE
	go mod tidy
	go mod vendor

.PHONY: FORCE
