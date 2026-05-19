# Web Server Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY common/ ./common/
COPY proto/ ./proto/
COPY logic-server/ ./logic-server/
COPY web-server/ ./web-server/

# 构建
WORKDIR /app/web-server
RUN CGO_ENABLED=0 GOOS=linux go build -o web-server .

# 运行时镜像 - 使用scratch最小镜像
FROM alpine:latest
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/web-server/web-server .

EXPOSE 8080

CMD ["./web-server"]
