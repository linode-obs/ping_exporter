ARG ARCH="amd64"
ARG OS="linux"
FROM golang:alpine3.18 AS builderimage
LABEL maintainer="Will Bollock <wbollock@gmail.com>"
WORKDIR /go/src/ping-exporter
COPY . .
RUN go build -o ping-exporter ../cmd/main.go

###################################################################

FROM golang:alpine3.18
COPY --from=builderimage /go/src/ping-exporter/ping-exporter /app/
WORKDIR /app

EXPOSE      9141
USER        nobody
ENTRYPOINT  [ "./ping-exporter" ]
