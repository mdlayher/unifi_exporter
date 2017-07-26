FROM alpine:latest

EXPOSE 9130
CMD ["/bin/unifi_exporter"]

RUN apk update ; apk add go ; apk add git ; apk add musl-dev ; \
    go get github.com/mdlayher/unifi_exporter/cmd/unifi_exporter; \
    mv ~/go/bin/unifi_exporter /bin/
