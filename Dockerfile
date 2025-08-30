# 多阶段构建
FROM golang:1.25-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的包
RUN apk add --no-cache git ca-certificates tzdata

# 复制 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o milonra-go main.go

# 最终镜像
FROM alpine:latest

# 安装必要的包
RUN apk --no-cache add ca-certificates tzdata

# 创建非特权用户
RUN addgroup -g 1000 milonra-go && \
    adduser -D -s /bin/sh -u 1000 -G milonra-go milonra-go

# 设置工作目录
WORKDIR /home/milonra-go

# 从构建阶段复制二进制文件
COPY --from=builder /app/milonra-go .

# 创建数据目录
RUN mkdir -p data logs && \
    chown -R milonra-go:milonra-go /home/milonra-go

# 切换到非特权用户
USER milonra-go

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./milonra-go"]
