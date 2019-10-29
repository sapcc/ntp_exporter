PKG    = github.com/sapcc/ntp_exporter
PREFIX = /usr

all: build/ntp_exporter

# NOTE: This repo uses Go modules, and uses a synthetic GOPATH at
# $(CURDIR)/.gopath that is only used for the build cache. $GOPATH/src/ is
# empty.
GO            = GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS =
GO_LDFLAGS    = -s -w

VERSION ?= $(shell sh util/find_version.sh)

build/ntp_exporter: *.go
	$(GO) install $(GO_BUILDFLAGS) -ldflags "$(GO_LDFLAGS) -X main.version=$(VERSION)" '$(PKG)'

install: build/ntp_exporter
	install -D -m 0755 build/ntp_exporter "$(DESTDIR)$(PREFIX)/bin/ntp_exporter"

vendor:
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: install vendor
