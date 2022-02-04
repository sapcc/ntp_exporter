FROM golang:1.17-alpine3.15 as builder
RUN apk add --no-cache make gcc musl-dev

COPY . /src
RUN make -C /src install PREFIX=/pkg GO_BUILDFLAGS='-mod vendor'

################################################################################

FROM alpine:3.15
MAINTAINER "Stefan Majewsky <stefan.majewsky@sap.com>"
LABEL source_repository="https://github.com/sapcc/ntp_exporter"

COPY --from=builder /pkg/ /usr/
ENTRYPOINT ["/usr/bin/ntp_exporter"]
