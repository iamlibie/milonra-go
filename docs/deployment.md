# 部署指南

本文档介绍如何部署 milonra-go 到不同环境。

## 🚀 快速部署

### 本地开发

```bash
git clone https://github.com/iamlibie/milonra-go.git
cd milonra-go
go mod tidy
go run main.go
```

### 编译部署

```bash
# 编译
go build -o milonra-go main.go

# 运行
./milonra-go
```

## 🐳 Docker 部署

### Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o milonra-go main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/milonra-go .
CMD ["./milonra-go"]
```

### docker-compose.yml

```yaml
version: '3.8'
services:
  milonra-go:
    build: .
    ports:
      - "8080:8080"
    environment:
      - BOT_ID=123456789
    volumes:
      - ./data:/root/data
    restart: unless-stopped
```

## ⚙️ 配置说明

### 环境变量

- `BOT_ID`: 机器人QQ号
- `LAGRANGE_URL`: Lagrange服务地址
- `API_PORT`: API端口（默认8080）

### 配置文件

创建 `config.json`：

```json
{
  "bot_id": 123456789,
  "lagrange_url": "ws://localhost:8080",
  "plugins": {
    "enabled": ["echo", "time", "userinfo"],
    "disabled": ["ai"]
  }
}
```

## 🔧 生产部署

### 系统服务 (systemd)

创建 `/etc/systemd/system/milonra-go.service`：

```ini
[Unit]
Description=milonra-go Service
After=network.target

[Service]
Type=simple
User=milonra-go
WorkingDirectory=/opt/milonra-go
ExecStart=/opt/milonra-go/milonra-go
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl enable milonra-go
sudo systemctl start milonra-go
```

## 🔒 安全配置

### 反向代理 (Nginx)

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### HTTPS 配置

```bash
# 使用 Let's Encrypt
sudo certbot --nginx -d your-domain.com
```

## 📊 监控与日志

### 日志配置

```go
// 在main.go中添加
log.SetOutput(&lumberjack.Logger{
    Filename:   "/var/log/milonra-go/app.log",
    MaxSize:    500, // megabytes
    MaxBackups: 3,
    MaxAge:     28,   // days
    Compress:   true,
})
```

### 健康检查

```bash
# 添加健康检查端点
curl http://localhost:8080/health
```

---

🚀 部署完成！更多配置选项请参考 [配置文档](configuration.md)。
