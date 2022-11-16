ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine as builder

RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY events events
COPY repository repository
COPY database database
COPY search search
COPY models models
COPY services/feed-service feed-service
COPY services/query-service query-service

RUN go install ./...

FROM alpine:3.11 as runner
WORKDIR /usr/bin/

COPY --from=builder /go/bin .
