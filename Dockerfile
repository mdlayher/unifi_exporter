FROM alpine
COPY ./unifi_exporter /bin/

USER nobody
EXPOSE 9130
ENTRYPOINT ["/bin/unifi_exporter"]
