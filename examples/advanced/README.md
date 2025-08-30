# 高级机器人示例

这个示例展示了如何构建一个功能更完整的机器人，包含：

## 功能特性

- 🌤️ **天气查询**: 查询城市天气信息
- 🧮 **计算器**: 数学表达式计算
- ⏰ **定时提醒**: 设置定时提醒功能  
- 🛡️ **管理功能**: 群管理员专用功能
- 📊 **数据统计**: 消息统计和用户活跃度
- 🎮 **小游戏**: 简单的互动游戏

## 使用方法

1. 复制此目录到新位置
2. 修改 `main.go` 中的机器人QQ号
3. 在 `plugins.go` 中添加管理员QQ号
4. 运行 `go mod tidy && go run main.go`

## 命令列表

- `天气 城市名` - 查询天气
- `计算 表达式` - 数学计算  
- `提醒 内容` - 设置提醒
- `状态` - 群状态查询（管理员）

## 扩展插件

在 `plugins/` 目录中添加更多插件：

```go
package plugins

func MyNewPlugin(bot plugin.Bot, e *event.MessageEvent) string {
    // 插件逻辑
    return "回复内容"
}

func init() {
    plugin.Register("my-plugin", MyNewPlugin)
}
```

## 配置示例

创建 `config.json`：

```json
{
    "bot_id": 123456789,
    "admin_users": [123456789, 987654321],
    "weather_api_key": "your-api-key",
    "database_url": "sqlite:./data.db"
}
```

## 部署说明

参考主项目的 [部署指南](../../docs/deployment.md)。