package plugin

import (
	"fmt"

	"github.com/iamlibie/milonra-go/event"
)

// Bot 接口，避免循环导入
type Bot interface {
	WriteJSON(v interface{}) error // 线程安全地写入JSON数据
	GetSelfID() int64              // 返回机器人ID
}

// PluginFunc 插件函数类型：输入bot实例和消息事件，输出回复
type PluginFunc func(bot Bot, event *event.MessageEvent) string

// 存储所有注册的插件
var plugins = make(map[string]PluginFunc)

// Register 注册一个插件
func Register(name string, fn PluginFunc) {
	plugins[name] = fn
	fmt.Printf("插件已注册: %s\n", name)
}

func GetPlugins() map[string]PluginFunc {
	return plugins
}
