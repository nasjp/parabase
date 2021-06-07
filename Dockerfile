FROM golang:1.16.5


WORKDIR /app

COPY go.mod /app/go.mod
COPY go.sum /app/go.sum

ENV IS_DOCKER true

RUN go mod download
