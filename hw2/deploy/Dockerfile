## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY main.go ./

RUN go build -o /ping main.go

## Deploy
FROM ubuntu:20.04

RUN apt-get update && yes | apt-get upgrade
RUN yes | apt-get install iputils-ping

WORKDIR /

COPY --from=build /ping /ping

ENTRYPOINT ["/ping"]