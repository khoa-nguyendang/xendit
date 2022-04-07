FROM golang:alpine as builder
COPY . /xendit-server
WORKDIR  /xendit-server
RUN apk add git && apk add build-base
ENV CGO_ENABLED 1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go install ./cmd/... 
