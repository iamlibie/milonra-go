// 简单机器人示例
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iamlibie/milonra-go/bot"
	"github.com/iamlibie/milonra-go/event"
	"github.com/iamlibie/milonra-go/plugin"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// 自定义插件示例
func helloPlugin(bot plugin.Bot, e *event.MessageEvent) string {
	if e.Message == "hello" {
		return fmt.Sprintf("Hello, %d! 欢迎使用 milonra-go!", e.UserID)
	}
	return ""
}

func pingPlugin(bot plugin.Bot, e *event.MessageEvent) string {
	if e.Message == "ping" {
		return "pong! 🏓"
	}
	return ""
}

func main() {
	// 注册自定义插件
	plugin.Register("hello", helloPlugin)
	plugin.Register("ping", pingPlugin)

	// 设置WebSocket处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("WebSocket升级失败:", err)
			return
		}
		defer conn.Close()

		fmt.Println("✅ 机器人已连接！")

		// 创建机器人实例
		botInstance := &bot.Bot{
			Conn:   conn,
			SelfID: 123456789, // 请替换为您的机器人QQ号
		}

		// 消息处理循环
		for {
			var data map[string]interface{}
			err := conn.ReadJSON(&data)
			if err != nil {
				fmt.Println("❌ 读取消息失败:", err)
				break
			}

			// 交给Bot处理
			botInstance.HandleMessage(data)
		}
	})

	fmt.Println("🚀 简单机器人启动，监听 :8080")
	fmt.Println("📝 支持的命令: hello, ping")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
