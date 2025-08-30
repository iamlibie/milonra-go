// 自定义插件示例
package myplugins

import (
	"fmt"
	"strings"
	"time"

	"github.com/iamlibie/milonra-go/api"
	"github.com/iamlibie/milonra-go/event"
	"github.com/iamlibie/milonra-go/plugin"
)

// 天气查询插件示例
func WeatherPlugin(bot plugin.Bot, e *event.MessageEvent) string {
	if strings.HasPrefix(e.Message, "天气 ") {
		city := strings.TrimPrefix(e.Message, "天气 ")
		// 这里应该调用实际的天气API
		return fmt.Sprintf("📍 %s 的天气：☀️ 晴天 25°C", city)
	}
	return ""
}

// 计算器插件示例
func CalculatorPlugin(bot plugin.Bot, e *event.MessageEvent) string {
	if strings.HasPrefix(e.Message, "计算 ") {
		expression := strings.TrimPrefix(e.Message, "计算 ")
		// 简单的加法示例
		if strings.Contains(expression, "+") {
			parts := strings.Split(expression, "+")
			if len(parts) == 2 {
				return fmt.Sprintf("💡 %s = 结果（请实现实际计算逻辑）", expression)
			}
		}
		return "❌ 暂只支持简单加法，格式：计算 1+2"
	}
	return ""
}

// 定时提醒插件示例
func ReminderPlugin(bot plugin.Bot, e *event.MessageEvent) string {
	if strings.HasPrefix(e.Message, "提醒 ") {
		reminder := strings.TrimPrefix(e.Message, "提醒 ")

		// 启动异步定时任务（示例：5秒后提醒）
		go func() {
			time.Sleep(5 * time.Second)
			msg := fmt.Sprintf("⏰ 提醒：%s", reminder)

			if e.GroupID != 0 {
				api.SendGroupMessage(bot, e.GroupID, msg)
			} else {
				api.SendPrivateMessage(bot, e.UserID, msg)
			}
		}()

		return "✅ 已设置提醒，将在5秒后提醒您"
	}
	return ""
}

// 管理员插件示例
func AdminPlugin(bot plugin.Bot, e *event.MessageEvent) string {
	// 定义管理员QQ号列表
	adminUsers := map[int64]bool{
		123456789: true, // 请替换为实际的管理员QQ号
	}

	if !adminUsers[e.UserID] {
		return ""
	}

	if e.Message == "状态" && e.GroupID != 0 {
		// 获取群信息
		groupInfo, err := api.GetGroupInfo(bot, e.GroupID)
		if err != nil {
			return "❌ 获取群信息失败"
		}

		return fmt.Sprintf("📊 群状态：\n群名：%s\n成员数：%d",
			groupInfo.GroupName, groupInfo.MemberCount)
	}

	return ""
}

func init() {
	plugin.Register("weather", WeatherPlugin)
	plugin.Register("calculator", CalculatorPlugin)
	plugin.Register("reminder", ReminderPlugin)
	plugin.Register("admin", AdminPlugin)
}
