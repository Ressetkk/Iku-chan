FROM golang:alpine as builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ikud main.go

FROM alpine:3.12 as certs
RUN apk add -U --no-cache ca-certificates

FROM scratch

WORKDIR /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/bin/ikud /ikud

ENTRYPOINT ["/ikud"]