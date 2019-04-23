FROM golang:1.12-stretch AS builder

WORKDIR /workspace

COPY ./ ./
RUN go mod download

RUN make build/docker

FROM alpine:3.7

WORKDIR /opt

COPY --from=builder /data-collector /data-collector

RUN apk update && apk upgrade && apk add bash

ENTRYPOINT ["/data-collector"]
CMD ["/bin/bash"]