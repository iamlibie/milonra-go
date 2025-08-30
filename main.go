// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/iamlibie/milonra-go/bot"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// 设置Lagrange HTTP API端点（根据实际配置修改）

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("升级失败:", err)
			return
		}
		defer conn.Close()

		fmt.Println("✅ lagrange.OneBot 已连接！")

		// 创建机器人实例
		botInstance := &bot.Bot{
			Conn:   conn,
			SelfID: 1707899218,
		}

		// 消息处理循环
		for {
			var data map[string]interface{}
			err := conn.ReadJSON(&data)
			if err != nil {
				fmt.Println("❌ 读取消息失败:", err)
				break
			}

			// 交给 Bot 处理
			botInstance.HandleMessage(data)
		}
	})

	fmt.Println("🚀 Go 机器人启动，监听 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
