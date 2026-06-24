# DeepL MCP Server (Go)
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /deepl-mcp-server ./cmd

FROM alpine:3.21

RUN apk add --no-cache ca-certificates curl

WORKDIR /app

COPY --from=builder /deepl-mcp-server .

EXPOSE 8080

ENV MCP_TRANSPORT=http
ENV MCP_HTTP_ADDR=:8080

# DeepL / AWS S3 对象存储配置均通过 env_file 或 docker run -e 注入，见 .env.example
# APP_OSS_* 使用 aws-sdk-go-v2/service/s3，endpoint 填你的 OSS/S3 服务地址

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -fsS http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["./deepl-mcp-server"]
