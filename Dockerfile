FROM golang:alpine as builder

WORKDIR /go/src/git.sbruder.de/simon/alertmanager-gotify/

COPY alertmanager-gotify.go .

RUN apk add --no-cache git upx ca-certificates

RUN go get -v \
    && CGO_ENABLED=0 go build -v -ldflags="-s -w" \
    && upx --ultra-brute alertmanager-gotify

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /go/src/git.sbruder.de/simon/alertmanager-gotify/alertmanager-gotify /alertmanager-gotify

USER 1000

ENTRYPOINT ["/alertmanager-gotify"]

EXPOSE 8081
