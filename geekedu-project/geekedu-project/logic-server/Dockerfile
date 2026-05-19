# Logic Server Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY common/ ./common/
COPY proto/ ./proto/
COPY logic-server/ ./logic-server/

# 构建
WORKDIR /app/logic-server
RUN CGO_ENABLED=0 GOOS=linux go build -o logic-server .

# 运行时镜像 - 使用scratch最小镜像
FROM alpine:latest
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/logic-server/logic-server .

EXPOSE 50051

CMD ["./logic-server"]
