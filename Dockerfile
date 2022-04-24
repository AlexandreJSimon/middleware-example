FROM golang:1.17-stretch as builder_base

WORKDIR /middleware-example

ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

COPY go.mod .
COPY go.sum .

RUN go mod download

# STAGE 1: Build binaries
FROM builder_base as builder

COPY . /middleware-example

WORKDIR /middleware-example

RUN go build -a -installsuffix cgo -o publisher github.com/AlexandreJSimon/middleware-example/cmd/publisher
RUN go build -a -installsuffix cgo -o subscriber github.com/AlexandreJSimon/middleware-example/cmd/subscriber

# STAGE 2: Build server
FROM alpine as publisher

COPY --from=builder /middleware-example/publisher /go/bin/publisher

# STAGE 3: Build worker
FROM alpine as subscriber

COPY --from=builder /middleware-example/subscriber /go/bin/subscriber

ENTRYPOINT /go/bin/subscriber