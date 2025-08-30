# MilonraGo 

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![WebSocket](https://img.shields.io/badge/WebSocket-gorilla%2Fwebsocket-orange.svg)](https://github.com/gorilla/websocket)

一个基于 Go 语言开发的高性能、模块化机器人框架，支持 OneBot 标准协议。MilonraGo 旨在为开发者提供一个简单易用、功能强大的机器人开发平台，支持自定义配置和插件扩展。如想获得最新通知或是加快BUG反馈进度请加Q群881391061

## 特性

- **高性能**: 基于 Go 语言，支持高并发处理
- **插件化**: 灵活的插件系统，支持热插拔
- **标准协议**: 完整支持 OneBot 11 标准
- **实时通信**: 基于 WebSocket 的稳定连接
- **易于扩展**: 简洁的 API 设计，便于二次开发
- **完整文档**: 详细的开发文档和示例

## 快速开始

### 安装

```bash
go get github.com/iamlibie/milonra-go
```

### 快速使用 - SDK方式 

```go
package main

import (
    "github.com/iamlibie/milonra-go/sdk"
    "github.com/iamlibie/milonra-go/event"
)

func main() {
    // 创建配置
    config := &sdk.MiloraBotConfig{
        Port:      ":8080",     // 自定义端口
        BotID:     123456789,   // 你的机器人QQ号
        EnableLog: true,        // 启用日志
    }

    // 创建机器人实例
    mb := sdk.NewMiloraBot(config)

    // 注册插件
    mb.RegisterPlugin("hello", func(bot sdk.Bot, e *event.MessageEvent) string {
        if e.Message == "hello" {
            return "Hello! 欢迎使用 MilonraGo! 🚀"
        }
        return ""
    })

    // 启动服务
    mb.Start()
}
```

### 作为库使用

在您的项目中使用 MilonraGo：

```go
go mod init your-bot-project
go get github.com/iamlibie/milonra-go
```

然后创建您的机器人：

```go
package main

import (
    "github.com/iamlibie/milonra-go/sdk"
    // 可选：导入您的插件包
    // _ "./plugins"
)

func main() {
    // 使用SDK创建机器人
    config := sdk.DefaultConfig() // 使用默认配置
    config.BotID = 123456789      // 设置您的机器人QQ号
    config.Port = ":8080"         // 自定义端口
    config.PluginDir = "./my-plugins"  // 指定插件目录
    config.AutoLoadPlugins = true      // 启用自动加载插件

    mb := sdk.NewMiloraBot(config)
    mb.Start()
}
```

## 📖 文档

- [插件开发指南](docs/plugin-development.md)
- [部署指南](docs/deployment.md)
- [示例项目](examples/)
- [API 文档](https://pkg.go.dev/github.com/iamlibie/milonra-go)

## ⚙️ SDK 配置说明

### 配置选项

```go
import "time"

config := &sdk.MiloraBotConfig{
    Port:              ":8080",                    // 监听端口
    Host:              "",                         // 监听主机（空表示所有接口）
    BotID:             123456789,                  // 机器人QQ号
    ReadTimeout:       15 * time.Second,           // 读取超时
    WriteTimeout:      15 * time.Second,           // 写入超时
    EnableLog:         true,                       // 启用日志
    LogLevel:          "info",                     // 日志级别
    PluginDir:         "./plugins",                // 插件目录
    AutoLoadPlugins:   true,                       // 自动加载插件
    PluginFilePattern: "*.so",                     // 插件文件匹配模式
}
```

### 常用方法

```go
// 创建机器人实例
mb := sdk.NewMiloraBot(config)

// 设置配置
mb.SetPort(":9000")                    // 修改端口
mb.SetBotID(987654321)                // 修改机器人ID
mb.EnableLogging(false)               // 禁用日志
mb.SetLogLevel("debug")               // 设置日志级别

// 插件管理
mb.SetPluginDir("./custom-plugins")   // 设置插件目录
mb.SetAutoLoadPlugins(true)           // 启用自动加载插件
mb.LoadPluginsFromDir("./plugins")    // 手动加载指定目录的插件
mb.RegisterPlugin("name", pluginFunc) // 手动注册插件

// 启动和停止
mb.Start()                            // 启动服务
mb.Stop(10 * time.Second)             // 优雅停止

// 获取信息
count := mb.GetPluginCount()          // 插件数量
plugins := mb.ListPlugins()           // 插件列表
bot := mb.GetBot()                    // 获取Bot实例
```

### 监控端点

- **健康检查**: `http://localhost:8080/health`
- **状态信息**: `http://localhost:8080/status`

## 🧪 测试功能

项目内置了全功能测试插件，可以验证所有功能模块：

### 测试账号

- **机器人**: 1707899218 (消息收发)
- **用户1**: 1234567890 (普通用户测试)
- **用户2**: 9876543210 (多用户测试)
- **管理员**: 1111111111 (权限测试)
- **测试群**: 123456789

### 测试命令

```
测试开始    - 启动全功能测试
测试账号    - 查看所有测试账号
测试消息 [X] - 测试消息处理
测试API     - 测试API功能
测试插件    - 测试插件系统
测试事件    - 测试事件处理
测试性能    - 性能测试
帮助测试    - 显示测试帮助
测试结束    - 结束测试
```

## 🔌 插件开发

### 创建插件

插件是 milonra-go 的核心功能扩展方式。创建一个插件非常简单：

```go
package plugins

import (
    "github.com/iamlibie/milonra-go/plugin"
    "github.com/iamlibie/milonra-go/event"
)

func MyAwesomePlugin(bot plugin.Bot, e *event.MessageEvent) string {
    // 处理消息逻辑
    if e.Message == "ping" {
        return "pong!"
    }
    return "" // 返回空字符串表示不处理此消息
}

func init() {
    plugin.Register("awesome", MyAwesomePlugin)
}
```

### 内置插件示例

- **Echo插件**: 复读机功能，以!开头的消息会被复读
- **Time插件**: 时间查询功能
- **UserInfo插件**: 用户信息查询
- **AI插件**: AI对话功能（需要配置API）

更多插件示例请查看 [plugins目录](plugins/) 和 [examples目录](examples/)。

## 🏗️ 项目结构

```
milonra-go/
├── api/           # OneBot API 封装
├── bot/           # 机器人核心逻辑
├── event/         # 事件定义
├── plugin/        # 插件管理器
├── sdk/           # SDK用户接口
├── examples/      # 使用示例
├── docs/          # 文档
├── tests/         # 集成测试
└── main.go        # 示例启动文件
```

## 🚀 部署

### 编译运行

```bash
# 克隆项目
git clone https://github.com/iamlibie/milonra-go.git
cd milonra-go

# 安装依赖
go mod tidy

# 编译
go build -o milonra-go main.go

# 运行
./milonra-go
```

### Docker 部署

```bash
# 使用 Docker
docker build -t milonra-go .
docker run -p 8080:8080 milonra-go

# 或使用 docker-compose
docker-compose up -d
```

## 🤝 贡献


1. Fork 本仓库
2. 创建您的功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 📄 开源协议

本项目采用 MIT 协议 - 查看 [LICENSE](LICENSE) 文件了解详细信息。

## 🙏 致谢

- [OneBot](https://github.com/howmanybots/onebot) - 标准协议支持
- [Lagrange.OneBot](https://github.com/LagrangeDev/Lagrange.Core) - 无头QQ，与之通信
- [gorilla/websocket](https://github.com/gorilla/websocket) - WebSocket 实现
- 所有贡献者和用户

## 📞 联系我们

- 提交 Issue: [GitHub Issues](https://github.com/iamlibie/milonra-go/issues)
- 讨论: [QQGroup](https://qm.qq.com/q/OqjgLXweWK)

---

⭐ 不过是一个高中生想练习GO而糊出来的小框架，如果这个项目对您有帮助，请给我们一个星标！
