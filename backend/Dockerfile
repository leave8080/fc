# 使用多阶段构建
FROM golang:1.21 AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用轻量级基础镜像
FROM alpine:latest

# 安装ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 创建存储目录
RUN mkdir -p /root/crawl_results

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]