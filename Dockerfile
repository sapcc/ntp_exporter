FROM alpine:latest
MAINTAINER "Stefan Majewsky <stefan.majewsky@sap.com>"

ADD ntp_exporter /bin/ntp_exporter
ENTRYPOINT ["/bin/ntp_exporter"]
