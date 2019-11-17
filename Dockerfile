FROM golang:1.13.1-alpine3.10

LABEL maintainer="Julian Sauer julian_sauer@gmx.net"

COPY ./ /go/src/github.com/JulianSauer/RemoteController

WORKDIR /go/src/github.com/JulianSauer/RemoteController/
RUN apk add git \
 && go get -v ./... \
 && go build main.go \
 && apk del git

RUN adduser -D remote
USER remote

EXPOSE 8080

ENTRYPOINT ["./main"]
