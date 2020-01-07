FROM golang:1.13.5 as build-env

WORKDIR /go/src/adlacarte-history
ADD . /go/src/adlacarte-history

RUN go get -d -v ./...

RUN go build -o /go/bin/adlacarte-history

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/adlacarte-history /

EXPOSE 8080

CMD ["/adlacarte-history"]
ENTRYPOINT ["/adlacarte-history"]
