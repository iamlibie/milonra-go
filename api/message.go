package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MessageSegment 消息段
type MessageSegment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Message 消息构造器
type Message struct {
	segments []MessageSegment
}

// NewMessage 创建新的消息构造器
func NewMessage() *Message {
	return &Message{
		segments: make([]MessageSegment, 0),
	}
}

// Text 添加文本消息
func (m *Message) Text(text string) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "text",
		Data: map[string]interface{}{
			"text": text,
		},
	})
	return m
}

// At 添加@消息
func (m *Message) At(qq int64) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": fmt.Sprintf("%d", qq),
		},
	})
	return m
}

// AtAll @全体成员
func (m *Message) AtAll() *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "at",
		Data: map[string]interface{}{
			"qq": "all",
		},
	})
	return m
}

// Face 添加QQ表情
func (m *Message) Face(id int) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "face",
		Data: map[string]interface{}{
			"id": fmt.Sprintf("%d", id),
		},
	})
	return m
}

// Image 添加图片
func (m *Message) Image(file string, url ...string) *Message {
	data := map[string]interface{}{
		"file": file,
	}
	if len(url) > 0 && url[0] != "" {
		data["url"] = url[0]
	}
	m.segments = append(m.segments, MessageSegment{
		Type: "image",
		Data: data,
	})
	return m
}

// Record 添加语音
func (m *Message) Record(file string, url ...string) *Message {
	data := map[string]interface{}{
		"file": file,
	}
	if len(url) > 0 && url[0] != "" {
		data["url"] = url[0]
	}
	m.segments = append(m.segments, MessageSegment{
		Type: "record",
		Data: data,
	})
	return m
}

// Video 添加视频
func (m *Message) Video(file string, url ...string) *Message {
	data := map[string]interface{}{
		"file": file,
	}
	if len(url) > 0 && url[0] != "" {
		data["url"] = url[0]
	}
	m.segments = append(m.segments, MessageSegment{
		Type: "video",
		Data: data,
	})
	return m
}

// Music 添加音乐分享
func (m *Message) Music(typ string, id int64) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "music",
		Data: map[string]interface{}{
			"type": typ,
			"id":   fmt.Sprintf("%d", id),
		},
	})
	return m
}

// CustomMusic 添加自定义音乐分享
func (m *Message) CustomMusic(url, audio, title string, content ...string) *Message {
	data := map[string]interface{}{
		"type":  "custom",
		"url":   url,
		"audio": audio,
		"title": title,
	}
	if len(content) > 0 {
		data["content"] = content[0]
	}
	m.segments = append(m.segments, MessageSegment{
		Type: "music",
		Data: data,
	})
	return m
}

// Reply 添加回复
func (m *Message) Reply(messageID int32) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "reply",
		Data: map[string]interface{}{
			"id": fmt.Sprintf("%d", messageID),
		},
	})
	return m
}

// Forward 添加合并转发
func (m *Message) Forward(id string) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "forward",
		Data: map[string]interface{}{
			"id": id,
		},
	})
	return m
}

// Node 添加合并转发节点
func (m *Message) Node(userID int64, nickname string, content interface{}) *Message {
	var contentData interface{}

	// 支持字符串或Message对象作为内容
	switch v := content.(type) {
	case string:
		contentData = v
	case *Message:
		contentData = v.Build()
	default:
		contentData = fmt.Sprintf("%v", content)
	}

	m.segments = append(m.segments, MessageSegment{
		Type: "node",
		Data: map[string]interface{}{
			"user_id":  fmt.Sprintf("%d", userID),
			"nickname": nickname,
			"content":  contentData,
		},
	})
	return m
}

// NodeID 添加合并转发节点(通过消息ID)
func (m *Message) NodeID(messageID int32) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "node",
		Data: map[string]interface{}{
			"id": fmt.Sprintf("%d", messageID),
		},
	})
	return m
}

// Poke 戳一戳
func (m *Message) Poke(qq int64) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "poke",
		Data: map[string]interface{}{
			"qq": fmt.Sprintf("%d", qq),
		},
	})
	return m
}

// JSON 添加JSON消息(卡片消息)
func (m *Message) JSON(jsonStr string) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "json",
		Data: map[string]interface{}{
			"data": jsonStr,
		},
	})
	return m
}

// XML 添加XML消息
func (m *Message) XML(xmlStr string) *Message {
	m.segments = append(m.segments, MessageSegment{
		Type: "xml",
		Data: map[string]interface{}{
			"data": xmlStr,
		},
	})
	return m
}

// Build 构建消息为数组格式
func (m *Message) Build() []MessageSegment {
	return m.segments
}

// BuildJSON 构建消息为JSON字符串
func (m *Message) BuildJSON() (string, error) {
	data, err := json.Marshal(m.segments)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToCQCode 转换为CQ码字符串
func (m *Message) ToCQCode() string {
	var parts []string
	for _, seg := range m.segments {
		parts = append(parts, segmentToCQCode(seg))
	}
	return strings.Join(parts, "")
}

// segmentToCQCode 将消息段转换为CQ码
func segmentToCQCode(seg MessageSegment) string {
	if seg.Type == "text" {
		if text, ok := seg.Data["text"].(string); ok {
			return text
		}
		return ""
	}

	var params []string
	for k, v := range seg.Data {
		params = append(params, fmt.Sprintf("%s=%v", k, v))
	}

	if len(params) > 0 {
		return fmt.Sprintf("[CQ:%s,%s]", seg.Type, strings.Join(params, ","))
	}
	return fmt.Sprintf("[CQ:%s]", seg.Type)
}

// ParseCQCode 解析CQ码字符串为Message对象
func ParseCQCode(cqCode string) *Message {
	msg := NewMessage()

	// 简单的CQ码解析实现
	parts := strings.Split(cqCode, "[CQ:")
	for i, part := range parts {
		if i == 0 && part != "" {
			// 第一部分可能是纯文本
			msg.Text(part)
			continue
		}

		if part == "" {
			continue
		}

		// 查找CQ码结束位置
		endIdx := strings.Index(part, "]")
		if endIdx == -1 {
			// 不是有效的CQ码，当作文本处理
			msg.Text("[CQ:" + part)
			continue
		}

		// 解析CQ码内容
		cqContent := part[:endIdx]
		remaining := part[endIdx+1:]

		// 分离类型和参数
		typeParts := strings.SplitN(cqContent, ",", 2)
		cqType := typeParts[0]

		data := make(map[string]interface{})
		if len(typeParts) > 1 {
			// 解析参数
			params := strings.Split(typeParts[1], ",")
			for _, param := range params {
				kv := strings.SplitN(param, "=", 2)
				if len(kv) == 2 {
					data[kv[0]] = kv[1]
				}
			}
		}

		msg.segments = append(msg.segments, MessageSegment{
			Type: cqType,
			Data: data,
		})

		// 处理剩余的文本
		if remaining != "" {
			remainingParts := strings.Split(remaining, "[CQ:")
			for j, rPart := range remainingParts {
				if j == 0 && rPart != "" {
					msg.Text(rPart)
				} else if rPart != "" {
					// 递归处理剩余的CQ码
					parts = append(parts, rPart)
				}
			}
		}
	}

	return msg
}
