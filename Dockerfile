
FROM golang:1.20 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o mcp-server ./cmd/mcp-server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/mcp-server .
ENTRYPOINT ["./mcp-server"]
