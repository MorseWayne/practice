FROM golang:latest AS builder

# 设置工作目录
WORKDIR /app

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY main.go ./

# 构建应用（禁用 CGO 便于在精简镜像运行）
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags='-s -w' -o main .

# 使用轻量级的 alpine 镜像作为运行环境
FROM alpine:latest

# 安装 ca-certificates（用于 HTTPS）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
