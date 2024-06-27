FROM golang:1.22.2 AS builder
WORKDIR /go/src/app
COPY . .
RUN GOARCH=amd64 GOOS=linux go build -o spacelift-challenge .

FROM docker
COPY --from=builder /go/src/app/spacelift-challenge /usr/local/bin/spacelift-challenge
RUN apk add bash curl
