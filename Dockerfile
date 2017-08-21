FROM golang:alpine AS build
ADD . /go/src/github.com/mdlayher/unifi_exporter
WORKDIR /go/src/github.com/mdlayher/unifi_exporter
RUN apk --update add make

RUN make build

FROM alpine
WORKDIR /app
COPY --from=build /go/src/github.com/mdlayher/unifi_exporter /bin/

USER nobody
EXPOSE 9130
ENTRYPOINT ["/bin/unifi_exporter"]
