FROM golang:1.12-stretch AS builder

WORKDIR /workspace

COPY ./ ./
RUN go mod download

RUN make build/docker

FROM alpine:3.7

WORKDIR /opt

COPY --from=builder /web-service /web-service

ENTRYPOINT ["/web-service"]
CMD ["/bin/bash"]