all: ntp_exporter

# force people to use golangvend
GOCC := env GOPATH=$(CURDIR)/.gopath go
GOFLAGS := -ldflags '-s -w'

ntp_exporter: *.go
	$(GOCC) build $(GOFLAGS) -o $@ github.com/sapcc/ntp_exporter

vendor:
	@golangvend
.PHONY: vendor
