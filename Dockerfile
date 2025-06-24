
FROM golang:1.22 AS builder
WORKDIR /app
ENV GOPROXY=direct
ENV GOTOOLCHAIN=local
COPY . .
RUN go mod tidy && go build -o mcp-server ./cmd/mcp-server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mcp-server .
ENTRYPOINT ["./mcp-server"]
