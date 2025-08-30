package bot

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"

	"log"

	"github.com/iamlibie/milonra-go/api"
	"github.com/iamlibie/milonra-go/event"
	"github.com/iamlibie/milonra-go/plugin"
)

type Bot struct {
	Conn       *websocket.Conn
	SelfID     int64
	writeMutex sync.Mutex
}

// WriteJSON 线程安全地向WebSocket写入JSON数据
func (b *Bot) WriteJSON(v interface{}) error {
	b.writeMutex.Lock()
	defer b.writeMutex.Unlock()
	return b.Conn.WriteJSON(v)
}

// GetSelfID 实现api.BotAPI接口
func (b *Bot) GetSelfID() int64 {
	return b.SelfID
}

// SendMessage 发送群消息（封装 OneBot API）

// HandleMessage 处理收到的消息
func (b *Bot) HandleMessage(data map[string]interface{}) {
	// 首先检查是否为 API 响应
	api.HandleAPIResponse(data)

	// 检查是否为消息事件
	postType, ok := data["post_type"].(string)
	if !ok || postType != "message" {
		// 可能是 heartbeat、meta_event 或 API 响应等，直接返回
		return
	}

	// 检查消息类型（群聊或私聊）
	messageType, ok := data["message_type"].(string)
	if !ok || (messageType != "group" && messageType != "private") {
		// 不是群聊也不是私聊，直接返回
		return
	}

	// 提取用户ID（群聊和私聊都有）
	userID, ok := data["user_id"].(float64)
	if !ok {
		log.Printf("❌ user_id 不存在或不是 float64: %v", data["user_id"])
		return
	}

	// 提取时间戳
	time, ok := data["time"].(float64)
	if !ok {
		log.Printf("❌ time 不存在或不是 float64: %v", data["time"])
		return
	}

	// 创建消息事件
	msgEvent := &event.MessageEvent{
		UserID:     int64(userID),
		Message:    api.ExtractMessage(data),
		RawMessage: api.GetString(data, "raw_message"),
		Time:       int64(time),
		RawData:    data,
	}

	// 如果是群聊消息，设置群ID
	if messageType == "group" {
		groupID, ok := data["group_id"].(float64)
		if !ok {
			log.Printf("❌ group_id 不存在或不是 float64: %v", data["group_id"])
			return
		}
		msgEvent.GroupID = int64(groupID)
		fmt.Printf("[群:%d] 用户:%d 说: %s\n", msgEvent.GroupID, msgEvent.UserID, msgEvent.Message)
	} else {
		// 私聊消息
		fmt.Printf("[私聊] 用户:%d 说: %s\n", msgEvent.UserID, msgEvent.Message)
	}

	// 检查是否@了机器人
	if api.IsAtMe(msgEvent.Message, b.SelfID) {
		msgEvent.IsAtMe = true
	}

	// 调用各个插件处理
	for name, pluginFunc := range plugin.GetPlugins() {
		// 异步发送回复，避免阻塞
		go func(evt *event.MessageEvent, pf plugin.PluginFunc, name string) {
			reply := pf(b, evt)
			if reply == "" {
				return
			}
			log.Printf("插件%s被调用", name)
			var err error
			if evt.GroupID != 0 {
				// 群聊消息
				_, err = api.SendGroupMessage(b, evt.GroupID, reply)
			} else {
				// 私聊消息
				_, err = api.SendPrivateMessage(b, evt.UserID, reply)
			}
			if err != nil {
				log.Printf("❌ 发送消息失败: %v", err)
			}
		}(msgEvent, pluginFunc, name)
	}
}
