FROM alpine:latest

EXPOSE 9130
ENTRYPOINT ["/bin/unifi_exporter"]

ADD unifi_exporter /bin/unifi_exporter
