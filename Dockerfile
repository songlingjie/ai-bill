# 第一阶段：构建
FROM golang:1.22 AS builder

# 设置工作目录
WORKDIR /app

# 先复制 go.mod 和 go.sum（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 再复制源码
COPY . .

# 编译（生成静态二进制）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

# 第二阶段：运行
FROM alpine:latest

WORKDIR /root/

# 安装证书（避免 HTTPS 请求报错）
RUN apk add --no-cache ca-certificates

# 拷贝编译好的程序
COPY --from=builder /app/app .

# 暴露端口（按你项目改）
EXPOSE 8080

# 启动
CMD ["./app"]