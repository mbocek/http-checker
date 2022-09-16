# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

ENV CGO_ENABLED=0
RUN go build -o /http-checker

## Deploy
FROM gcr.io/distroless/base
WORKDIR /
COPY --from=build /http-checker /http-checker
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/http-checker"]
