package api

import (
	"encoding/json"
	"fmt"

	"github.com/iamlibie/milonra-go/plugin"
)

// SendGroupForwardMsg 发送群聊合并转发消息
func SendGroupForwardMsg(b plugin.Bot, groupID int64, messages []MessageSegment) (int32, error) {
	echo := generateEcho("send_group_forward_msg")
	params := map[string]interface{}{
		"group_id": groupID,
		"messages": messages,
	}
	data := map[string]interface{}{
		"action": "send_group_forward_msg",
		"params": params,
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return 0, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return 0, err
	}

	if resp.Status != "ok" {
		return 0, fmt.Errorf("发送群聊合并转发消息失败: %s", resp.Status)
	}

	var result struct {
		MessageID int32 `json:"message_id"`
	}
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return 0, err
	}

	return result.MessageID, nil
}

// SendPrivateForwardMsg 发送私聊合并转发消息
func SendPrivateForwardMsg(b plugin.Bot, userID int64, messages []MessageSegment) (int32, error) {
	echo := generateEcho("send_private_forward_msg")
	params := map[string]interface{}{
		"user_id":  userID,
		"messages": messages,
	}
	data := map[string]interface{}{
		"action": "send_private_forward_msg",
		"params": params,
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return 0, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return 0, err
	}

	if resp.Status != "ok" {
		return 0, fmt.Errorf("发送私聊合并转发消息失败: %s", resp.Status)
	}

	var result struct {
		MessageID int32 `json:"message_id"`
	}
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return 0, err
	}

	return result.MessageID, nil
}

// ForwardMessageBuilder 合并转发消息构建器
type ForwardMessageBuilder struct {
	nodes []MessageSegment
}

// NewForwardMessage 创建合并转发消息构建器
func NewForwardMessage() *ForwardMessageBuilder {
	return &ForwardMessageBuilder{
		nodes: make([]MessageSegment, 0),
	}
}

// AddNode 添加转发节点（自定义内容）
func (f *ForwardMessageBuilder) AddNode(userID int64, nickname string, content interface{}) *ForwardMessageBuilder {
	var contentData interface{}

	// 支持多种内容格式
	switch v := content.(type) {
	case string:
		// 纯文本
		contentData = []MessageSegment{{
			Type: "text",
			Data: map[string]interface{}{"text": v},
		}}
	case *Message:
		// Message对象
		contentData = v.Build()
	case []MessageSegment:
		// 消息段数组
		contentData = v
	default:
		// 其他格式转为字符串
		contentData = []MessageSegment{{
			Type: "text",
			Data: map[string]interface{}{"text": fmt.Sprintf("%v", content)},
		}}
	}

	node := MessageSegment{
		Type: "node",
		Data: map[string]interface{}{
			"user_id":  userID,
			"nickname": nickname,
			"content":  contentData,
		},
	}

	f.nodes = append(f.nodes, node)
	return f
}

// AddNodeByID 添加转发节点（引用消息ID）
func (f *ForwardMessageBuilder) AddNodeByID(messageID int32) *ForwardMessageBuilder {
	node := MessageSegment{
		Type: "node",
		Data: map[string]interface{}{
			"id": messageID,
		},
	}

	f.nodes = append(f.nodes, node)
	return f
}

// AddCustomNode 添加自定义时间和发送者的节点
func (f *ForwardMessageBuilder) AddCustomNode(userID int64, nickname string, content interface{}, time int64) *ForwardMessageBuilder {
	var contentData interface{}

	switch v := content.(type) {
	case string:
		contentData = []MessageSegment{{
			Type: "text",
			Data: map[string]interface{}{"text": v},
		}}
	case *Message:
		contentData = v.Build()
	case []MessageSegment:
		contentData = v
	default:
		contentData = []MessageSegment{{
			Type: "text",
			Data: map[string]interface{}{"text": fmt.Sprintf("%v", content)},
		}}
	}

	node := MessageSegment{
		Type: "node",
		Data: map[string]interface{}{
			"user_id":  userID,
			"nickname": nickname,
			"content":  contentData,
			"time":     time,
		},
	}

	f.nodes = append(f.nodes, node)
	return f
}

// Build 构建消息段数组
func (f *ForwardMessageBuilder) Build() []MessageSegment {
	return f.nodes
}

// SendToGroup 发送到群聊
func (f *ForwardMessageBuilder) SendToGroup(b plugin.Bot, groupID int64) (int32, error) {
	return SendGroupForwardMsg(b, groupID, f.nodes)
}

// SendToPrivate 发送到私聊
func (f *ForwardMessageBuilder) SendToPrivate(b plugin.Bot, userID int64) (int32, error) {
	return SendPrivateForwardMsg(b, userID, f.nodes)
}

// GetCustomFace 获取自定义表情
func GetCustomFace(b plugin.Bot) ([]string, error) {
	echo := generateEcho("fetch_custom_face")
	data := map[string]interface{}{
		"action": "fetch_custom_face",
		"params": map[string]interface{}{},
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return nil, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("获取自定义表情失败: %s", resp.Status)
	}

	var faces []string
	err = json.Unmarshal(resp.Data, &faces)
	if err != nil {
		return nil, err
	}

	return faces, nil
}

// GetMFaceKey 获取商城表情 key
func GetMFaceKey(b plugin.Bot, emojiIDs []string) ([]string, error) {
	echo := generateEcho("fetch_mface_key")
	data := map[string]interface{}{
		"action": "fetch_mface_key",
		"params": map[string]interface{}{
			"emoji_ids": emojiIDs,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return nil, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("获取商城表情key失败: %s", resp.Status)
	}

	var keys []string
	err = json.Unmarshal(resp.Data, &keys)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

// JoinGroupEmojiChain 加入群聊表情接龙
func JoinGroupEmojiChain(b plugin.Bot, groupID int64, messageID int32, emojiID int32) error {
	echo := generateEcho(".join_group_emoji_chain")
	data := map[string]interface{}{
		"action": ".join_group_emoji_chain",
		"params": map[string]interface{}{
			"group_id":   groupID,
			"message_id": messageID,
			"emoji_id":   emojiID,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("加入表情接龙失败: %s", resp.Status)
	}

	return nil
}

// JoinFriendEmojiChain 加入好友表情接龙
func JoinFriendEmojiChain(b plugin.Bot, userID int64, messageID int32, emojiID int32) error {
	echo := generateEcho(".join_friend_emoji_chain")
	data := map[string]interface{}{
		"action": ".join_friend_emoji_chain",
		"params": map[string]interface{}{
			"user_id":    userID,
			"message_id": messageID,
			"emoji_id":   emojiID,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("加入表情接龙失败: %s", resp.Status)
	}

	return nil
}

// AICharacter AI语音角色信息
type AICharacter struct {
	CharacterID   string `json:"character_id"`
	CharacterName string `json:"character_name"`
	PreviewURL    string `json:"preview_url"`
}

// AICharacterGroup AI语音角色分组
type AICharacterGroup struct {
	Type       string        `json:"type"`
	Characters []AICharacter `json:"characters"`
}

// GetAICharacters 获取群 AI 语音可用声色列表
func GetAICharacters(b plugin.Bot, groupID int64, chatType int) ([]AICharacterGroup, error) {
	echo := generateEcho("get_ai_characters")
	data := map[string]interface{}{
		"action": "get_ai_characters",
		"params": map[string]interface{}{
			"group_id":  groupID,
			"chat_type": chatType, // 1: 朗读, 2: 说唱
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return nil, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("获取AI语音列表失败: %s", resp.Status)
	}

	var characters []AICharacterGroup
	err = json.Unmarshal(resp.Data, &characters)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

// SendGroupAIVoice 发送群 AI 语音
func SendGroupAIVoice(b plugin.Bot, groupID int64, characterID string, text string) (int32, error) {
	echo := generateEcho("send_group_ai_voice")
	data := map[string]interface{}{
		"action": "send_group_ai_voice",
		"params": map[string]interface{}{
			"group_id":     groupID,
			"character_id": characterID,
			"text":         text,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return 0, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return 0, err
	}

	if resp.Status != "ok" {
		return 0, fmt.Errorf("发送AI语音失败: %s", resp.Status)
	}

	var result struct {
		MessageID int32 `json:"message_id"`
	}
	err = json.Unmarshal(resp.Data, &result)
	if err != nil {
		return 0, err
	}

	return result.MessageID, nil
}
