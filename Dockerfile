FROM golang:1.13.5 as build-env

WORKDIR /go/src/RemoteController
ADD . /go/src/RemoteController

RUN go get -d -v ./...

RUN go build -o /go/bin/RemoteController

FROM scratch
COPY --from=build-env /go/bin/RemoteController /

EXPOSE 8080

ENTRYPOINT ["./main"]
