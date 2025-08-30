# 插件开发指南

本文档将指导您如何为 milonra-go 开发自定义插件。

## 📚 目录

- [快速开始](#快速开始)
- [插件结构](#插件结构)
- [API 参考](#api-参考)
- [事件类型](#事件类型)
- [消息构建](#消息构建)
- [高级功能](#高级功能)
- [最佳实践](#最佳实践)
- [示例插件](#示例插件)

## 🚀 快速开始

### 创建第一个插件

创建一个新的 `.go` 文件，例如 `myplugin.go`：

```go
package plugins

import (
    "github.com/iamlibie/milonra-go/plugin"
    "github.com/iamlibie/milonra-go/event"
)

// 插件处理函数
func MyFirstPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    if e.Message == "hello" {
        return "Hello! 欢迎使用 milonra-go!"
    }
    return "" // 返回空字符串表示不处理此消息
}

// 插件注册
func init() {
    plugin.Register("my-first-plugin", MyFirstPlugin)
}
```

### 插件文件放置

将插件文件放在以下位置之一：

1. **项目内插件**: `plugins/` 目录下
2. **外部插件**: 任何 Go 包，通过 import 引入

## 🏗️ 插件结构

### 插件函数签名

所有插件必须遵循以下函数签名：

```go
func PluginName(bot plugin.Bot, e *event.MessageEvent) string
```

**参数说明**：
- `bot plugin.Bot`: 机器人接口，用于发送消息和调用API
- `e *event.MessageEvent`: 消息事件，包含消息内容和元数据
- 返回值 `string`: 要回复的消息内容，空字符串表示不回复

### 插件注册

在 `init()` 函数中注册插件：

```go
func init() {
    plugin.Register("plugin-name", PluginFunction)
}
```

## 🔧 API 参考

### Bot 接口

`plugin.Bot` 接口提供以下方法：

```go
type Bot interface {
    WriteJSON(v interface{}) error  // 发送JSON数据到WebSocket
    GetSelfID() int64              // 获取机器人自身的QQ号
}
```

### 常用API函数

```go
import "github.com/iamlibie/milonra-go/api"

// 发送群消息
api.SendGroupMessage(bot, groupID, message)

// 发送私聊消息
api.SendPrivateMessage(bot, userID, message)

// 获取群成员信息
api.GetGroupMemberInfo(bot, groupID, userID)

// 获取群信息
api.GetGroupInfo(bot, groupID)

// 撤回消息
api.DeleteMessage(bot, messageID)
```

## 📝 事件类型

### MessageEvent 结构

```go
type MessageEvent struct {
    GroupID    int64                  // 群号（私聊时为0）
    UserID     int64                  // 发送者QQ号
    Message    string                 // 消息内容（纯文本）
    RawMessage string                 // 原始消息（包含CQ码）
    Nickname   string                 // 发送者昵称
    Time       int64                  // 消息时间戳
    IsAtMe     bool                   // 是否@了机器人
    RawData    map[string]interface{} // 原始JSON数据
}
```

### 消息类型判断

```go
func MyPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    // 判断是群聊还是私聊
    if e.GroupID != 0 {
        // 群聊消息
        return "这是群聊消息"
    } else {
        // 私聊消息
        return "这是私聊消息"
    }
}
```

## 💬 消息构建

### 简单文本消息

```go
return "这是一条简单的文本消息"
```

### 复杂消息构建

使用 `api.Message` 构建复杂消息：

```go
import "github.com/iamlibie/milonra-go/api"

func ComplexMessagePlugin(bot plugin.Bot, e *event.MessageEvent) string {
    msg := api.NewMessage().
        At(e.UserID).                    // @用户
        Text(" 你好！").                   // 文本
        Image("path/to/image.jpg").      // 图片
        Face(1)                          // QQ表情

    return msg.ToCQCode() // 转换为CQ码字符串
}
```

### 支持的消息类型

- 📝 **文本**: `.Text("文本内容")`
- 👤 **@**: `.At(userID)` 或 `.AtAll()`
- 😀 **表情**: `.Face(faceID)`
- 🖼️ **图片**: `.Image("file_path")`
- 🎵 **音频**: `.Record("file_path")`
- 🎬 **视频**: `.Video("file_path")`
- 🎶 **音乐**: `.Music("qq", musicID)`
- 💾 **JSON卡片**: `.JSON(jsonString)`

## 🔧 高级功能

### 异步处理

```go
func AsyncPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    if e.Message == "async" {
        // 启动异步任务
        go func() {
            time.Sleep(5 * time.Second)
            api.SendGroupMessage(bot, e.GroupID, "异步任务完成！")
        }()
        return "异步任务已启动..."
    }
    return ""
}
```

### 状态管理

```go
var userState = make(map[int64]string)
var stateMutex sync.RWMutex

func StatefulPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    stateMutex.Lock()
    defer stateMutex.Unlock()

    state := userState[e.UserID]

    switch state {
    case "waiting_name":
        userState[e.UserID] = ""
        return fmt.Sprintf("你好，%s！", e.Message)
    default:
        if e.Message == "设置姓名" {
            userState[e.UserID] = "waiting_name"
            return "请输入你的姓名："
        }
    }
    return ""
}
```

### 数据持久化

```go
import (
    "encoding/json"
    "os"
)

type UserData struct {
    UserID int64  `json:"user_id"`
    Name   string `json:"name"`
    Level  int    `json:"level"`
}

func SaveUserData(data *UserData) error {
    file, err := os.Create(fmt.Sprintf("data/%d.json", data.UserID))
    if err != nil {
        return err
    }
    defer file.Close()

    return json.NewEncoder(file).Encode(data)
}

func LoadUserData(userID int64) (*UserData, error) {
    file, err := os.Open(fmt.Sprintf("data/%d.json", userID))
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var data UserData
    err = json.NewDecoder(file).Decode(&data)
    return &data, err
}
```

## ✅ 最佳实践

### 1. 错误处理

```go
func SafePlugin(bot plugin.Bot, e *event.MessageEvent) string {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("插件发生panic: %v", r)
        }
    }()

    // 插件逻辑...
    return "安全的插件响应"
}
```

### 2. 权限控制

```go
var adminUsers = map[int64]bool{
    123456789: true, // 管理员QQ号
}

func AdminOnlyPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    if !adminUsers[e.UserID] {
        return "权限不足"
    }

    // 管理员功能...
    return "管理员命令执行成功"
}
```

### 3. 命令解析

```go
import "strings"

func CommandPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    if !strings.HasPrefix(e.Message, "/") {
        return ""
    }

    parts := strings.Fields(e.Message)
    if len(parts) == 0 {
        return ""
    }

    command := parts[0][1:] // 去除 "/"
    args := parts[1:]

    switch command {
    case "help":
        return "可用命令: /help, /weather, /time"
    case "weather":
        if len(args) > 0 {
            return fmt.Sprintf("%s的天气信息...", args[0])
        }
        return "请提供城市名称"
    case "time":
        return time.Now().Format("2006-01-02 15:04:05")
    default:
        return "未知命令，使用 /help 查看帮助"
    }
}
```

### 4. 频率限制

```go
import (
    "sync"
    "time"
)

var (
    lastCall = make(map[int64]time.Time)
    callMutex sync.RWMutex
)

func RateLimitedPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    callMutex.Lock()
    defer callMutex.Unlock()

    if last, exists := lastCall[e.UserID]; exists {
        if time.Since(last) < 10*time.Second {
            return "请稍后再试（10秒冷却）"
        }
    }

    lastCall[e.UserID] = time.Now()
    return "命令执行成功"
}
```

## 📚 示例插件

查看 `plugins/` 目录下的示例插件：

- **echo.go**: 回声插件，复读消息
- **time.go**: 时间查询插件
- **userinfo.go**: 用户信息查询插件
- **ai_li.go**: AI对话插件

更多高级示例请查看 `examples/` 目录。

## 🔍 调试技巧

### 1. 日志记录

```go
import "log"

func DebugPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    log.Printf("收到消息: %s，来自用户: %d", e.Message, e.UserID)
    // 插件逻辑...
    return "调试信息已记录"
}
```

### 2. 消息内容检查

```go
func InspectPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    if e.Message == "debug" {
        return fmt.Sprintf("调试信息:\n群号: %d\n用户: %d\n消息: %s\n原始: %s",
            e.GroupID, e.UserID, e.Message, e.RawMessage)
    }
    return ""
}
```

## 🚀 发布插件

### 1. 创建独立包

```bash
mkdir my-plugin
cd my-plugin
go mod init github.com/iamlibie/my-plugin
```

### 2. 插件代码结构

```go
// my-plugin/plugin.go
package myplugin

import (
    "github.com/iamlibie/milonra-go/plugin"
    "github.com/iamlibie/milonra-go/event"
)

func MyAwesomeFeature(bot plugin.Bot, e *event.MessageEvent) string {
    // 插件逻辑
    return "Awesome!"
}

func Register() {
    plugin.Register("my-awesome-plugin", MyAwesomeFeature)
}
```

### 3. 使用外部插件

```go
// 在主项目中
import (
    _ "github.com/iamlibie/my-plugin"
)

func init() {
    myplugin.Register()
}
```

---

🎉 恭喜！您现在已经掌握了 milonra-go 插件开发的基础知识。开始创建您自己的插件吧！

如有问题，请查看 [示例代码](../examples/) 或 [提交Issue](https://github.com/iamlibie/milonra-go/issues)。
