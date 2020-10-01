FROM golang:alpine as builder

WORKDIR /src
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ikud main.go

FROM scratch

WORKDIR /
COPY --from=builder /src/bin/ikud /ikud

ENTRYPOINT ["/ikud"]