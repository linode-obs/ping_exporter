ARG ARCH="amd64"
ARG OS="linux"
FROM golang:1.27.3-alpine AS builderimage
LABEL maintainer="Will Bollock <wbollock@gmail.com>"
WORKDIR /go/src/ping-exporter
COPY . .
RUN go build -o ping-exporter cmd/main.go

###################################################################

FROM golang:1.27.3-alpine
COPY --from=builderimage /go/src/ping-exporter/ping-exporter /app/
WORKDIR /app

RUN apk add --no-cache libcap
RUN setcap cap_net_raw=+ep ./ping-exporter

EXPOSE      9141
USER        nobody
ENTRYPOINT  [ "./ping-exporter" ]
