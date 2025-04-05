FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

ARG VERSION=main
ARG MODULE=client
ARG ARGS=""

WORKDIR /app

RUN git clone --branch ${VERSION} https://github.com/ecromaneli/woole.git .

RUN cd cmd/${MODULE} && go build -o ../../bin/woole .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/woole .

ENTRYPOINT ["./woole"]