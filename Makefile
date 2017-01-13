all: ntp_exporter

# force people to use golangvend
GOCC := env GOPATH=$(CURDIR)/.gopath go
GOCCFLAGS :=
GOLDFLAGS := -s -w

VERSION ?= $(shell git describe --tags --dirty)

ntp_exporter: *.go
	$(GOCC) build $(GOCCFLAGS) -ldflags "$(GOLDFLAGS) -X main.version=$(VERSION)" -o $@ github.com/sapcc/ntp_exporter

vendor:
	@golangvend
.PHONY: vendor
