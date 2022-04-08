FROM golang:alpine as builder
COPY . /xendit-server
WORKDIR  /xendit-server
RUN apk add git && apk add build-base
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# RUN go build -tags=jsoniter -o /xendit-server/cmd/... .

# FROM scratch
# COPY --from=builder /xendit-server/cmd /cmd
# RUN go install /cmd
RUN go install  -tags="musl,jsoniter"  ./cmd/...