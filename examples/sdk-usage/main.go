// 使用MiloraBot SDK的简单示例
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamlibie/milonra-go/event"
	"github.com/iamlibie/milonra-go/sdk"
)

// 自定义插件示例
func helloPlugin(bot sdk.Bot, e *event.MessageEvent) string {
	if e.Message == "hello" {
		return "你好！我是使用SDK创建的机器人 🤖"
	}
	return ""
}

func pingPlugin(bot sdk.Bot, e *event.MessageEvent) string {
	if e.Message == "ping" {
		return "pong! 🏓"
	}
	return ""
}

func helpPlugin(bot sdk.Bot, e *event.MessageEvent) string {
	if e.Message == "help" || e.Message == "帮助" {
		return `📖 可用命令：
• hello - 打招呼
• ping - 测试连接
• help - 显示帮助
• status - 显示状态`
	}
	return ""
}

func statusPlugin(bot sdk.Bot, e *event.MessageEvent) string {
	if e.Message == "status" || e.Message == "状态" {
		return "✅ 机器人运行正常！"
	}
	return ""
}

func main() {
	// 创建配置
	config := &sdk.MiloraBotConfig{
		Port:      ":8080",   // 自定义端口
		BotID:     123456789, // 你的机器人QQ号
		EnableLog: true,      // 启用日志
		LogLevel:  "info",    // 日志级别
	}

	// 创建MiloraBot实例
	mb := sdk.NewMiloraBot(config)

	// 注册插件
	mb.RegisterPlugin("hello", helloPlugin)
	mb.RegisterPlugin("ping", pingPlugin)
	mb.RegisterPlugin("help", helpPlugin)
	mb.RegisterPlugin("status", statusPlugin)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务
	go func() {
		if err := mb.Start(); err != nil {
			log.Printf("❌ 服务启动失败: %v", err)
		}
	}()

	// 等待停止信号
	<-sigChan
	log.Println("🛑 收到停止信号，正在关闭服务...")

	// 优雅关闭
	if err := mb.Stop(10 * time.Second); err != nil {
		log.Printf("❌ 服务关闭失败: %v", err)
	}
}
