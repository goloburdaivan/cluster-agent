FROM golang:1.25.0-bookworm AS builder

RUN apt-get update && apt-get install -y clang llvm curl git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/internal/ebpf/headers/bpf
RUN curl -LO https://raw.githubusercontent.com/libbpf/libbpf/master/src/bpf_helpers.h \
    && curl -LO https://raw.githubusercontent.com/libbpf/libbpf/master/src/bpf_helper_defs.h \
    && curl -LO https://raw.githubusercontent.com/libbpf/libbpf/master/src/bpf_core_read.h \
    && curl -LO https://raw.githubusercontent.com/libbpf/libbpf/master/src/bpf_endian.h \
    && curl -LO https://raw.githubusercontent.com/libbpf/libbpf/master/src/bpf_tracing.h \
    && curl -LO https://raw.githubusercontent.com/libbpf/libbpf/master/src/bpf_common.h

WORKDIR /app

RUN go generate ./internal/ebpf/...

RUN CGO_ENABLED=0 GOOS=linux go build -o agent ./cmd/agent

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/agent .

CMD ["./agent"]