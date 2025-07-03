FROM golang:1.22 AS builder
WORKDIR /app
ENV GOPROXY=direct

# Download and unpack crictl
ENV CRICTL_VERSION=v1.28.0
RUN curl -LO https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz && \
    tar -xzf crictl-${CRICTL_VERSION}-linux-amd64.tar.gz && \
    mv crictl /usr/local/bin/ && \
    chmod +x /usr/local/bin/crictl && \
    rm crictl-${CRICTL_VERSION}-linux-amd64.tar.gz

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mcp-server ./cmd/mcp-server

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates curl iproute2 ethtool util-linux
COPY --from=builder /usr/local/bin/crictl /usr/local/bin/crictl
COPY --from=builder /app/mcp-server .
RUN chmod +x /usr/local/bin/crictl
ENTRYPOINT ["./mcp-server"]

