FROM alpine:latest

EXPOSE 9130

RUN apk add --update --virtual build-deps go git musl-dev && \
    go get github.com/mdlayher/unifi_exporter/cmd/unifi_exporter && \
    mv ~/go/bin/unifi_exporter /bin/ && \
    apk del build-deps && \
    rm -rf /var/cache/apk/* ~/go/

USER nobody
ENTRYPOINT ["/bin/unifi_exporter"]
