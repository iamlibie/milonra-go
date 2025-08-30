# SDK 使用示例

这个示例展示了如何使用 milonra-go SDK 创建自定义机器人。

## 📖 特性

- **🔧 可配置端口**: 自定义监听端口
- **🤖 自定义机器人ID**: 设置你的机器人QQ号
- **📝 日志控制**: 启用/禁用日志，设置日志级别
- **🔌 插件管理**: 轻松注册和管理插件
- **🛑 优雅关闭**: 支持信号处理和优雅关闭

## 🚀 快速开始

1. 修改 `main.go` 中的机器人QQ号：
```go
config := &sdk.MiloraBotConfig{
    BotID: 你的机器人QQ号, // 修改这里
    Port:  ":8080",        // 可以自定义端口
}
```

2. 运行机器人：
```bash
go run main.go
```

3. 连接你的OneBot客户端到 `ws://localhost:8080`

## 🔧 配置选项

```go
config := &sdk.MiloraBotConfig{
    Port:         ":8080",                    // 监听端口
    BotID:        123456789,                  // 机器人QQ号
    EnableLog:    true,                       // 启用日志
    LogLevel:     "info",                     // 日志级别
    ReadTimeout:  15 * time.Second,           // 读取超时
    WriteTimeout: 15 * time.Second,           // 写入超时
}
```

## 🔌 插件示例

```go
func myPlugin(bot sdk.Bot, e *event.MessageEvent) string {
    if e.Message == "test" {
        return "测试成功！"
    }
    return ""
}

// 注册插件
mb.RegisterPlugin("test", myPlugin)
```

## 📋 内置命令

- `hello` - 打招呼
- `ping` - 测试连接
- `help` - 显示帮助
- `status` - 显示状态

## 🔍 监控端点

- 健康检查: `http://localhost:8080/health`
- 状态信息: `http://localhost:8080/status`
