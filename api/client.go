package api

import (
	"encoding/json"
	"fmt"

	"github.com/iamlibie/milonra-go/plugin"
)

// SendGroupMessage 发送群消息
func SendGroupMessage(b plugin.Bot, groupID int64, message interface{}) (int32, error) {
	// 支持字符串、Message对象或MessageSegment数组
	var msg interface{}
	switch v := message.(type) {
	case string:
		msg = v
	case *Message:
		msg = v.Build()
	case []MessageSegment:
		msg = v
	default:
		msg = fmt.Sprintf("%v", message)
	}

	echo := generateEcho("send_group_msg")
	data := map[string]interface{}{
		"action": "send_group_msg",
		"params": map[string]interface{}{
			"group_id": groupID,
			"message":  msg,
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
		return 0, fmt.Errorf("发送失败: %s", resp.Status)
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

// SendPrivateMessage 发送私聊消息
func SendPrivateMessage(b plugin.Bot, userID int64, message interface{}) (int32, error) {
	var msg interface{}
	switch v := message.(type) {
	case string:
		msg = v
	case *Message:
		msg = v.Build()
	case []MessageSegment:
		msg = v
	default:
		msg = fmt.Sprintf("%v", message)
	}

	echo := generateEcho("send_private_msg")
	data := map[string]interface{}{
		"action": "send_private_msg",
		"params": map[string]interface{}{
			"user_id": userID,
			"message": msg,
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
		return 0, fmt.Errorf("发送失败: %s", resp.Status)
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

// SendMsg 通用发送消息
func SendMsg(b plugin.Bot, messageType string, userID, groupID int64, message interface{}) (int32, error) {
	var msg interface{}
	switch v := message.(type) {
	case string:
		msg = v
	case *Message:
		msg = v.Build()
	case []MessageSegment:
		msg = v
	default:
		msg = fmt.Sprintf("%v", message)
	}

	echo := generateEcho("send_msg")
	params := map[string]interface{}{
		"message": msg,
	}

	if messageType != "" {
		params["message_type"] = messageType
	}
	if userID != 0 {
		params["user_id"] = userID
	}
	if groupID != 0 {
		params["group_id"] = groupID
	}

	data := map[string]interface{}{
		"action": "send_msg",
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
		return 0, fmt.Errorf("发送失败: %s", resp.Status)
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

// DeleteMsg 撤回消息
func DeleteMsg(b plugin.Bot, messageID int32) error {
	echo := generateEcho("delete_msg")
	data := map[string]interface{}{
		"action": "delete_msg",
		"params": map[string]interface{}{
			"message_id": messageID,
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
		return fmt.Errorf("撤回失败: %s", resp.Status)
	}

	return nil
}

// GetMsg 获取消息
func GetMsg(b plugin.Bot, messageID int32) (*MessageInfo, error) {
	echo := generateEcho("get_msg")
	data := map[string]interface{}{
		"action": "get_msg",
		"params": map[string]interface{}{
			"message_id": messageID,
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
		return nil, fmt.Errorf("获取消息失败: %s", resp.Status)
	}

	var msgInfo MessageInfo
	err = json.Unmarshal(resp.Data, &msgInfo)
	if err != nil {
		return nil, err
	}

	return &msgInfo, nil
}

// GetForwardMsg 获取合并转发消息
func GetForwardMsg(b plugin.Bot, id string) (*ForwardMessage, error) {
	echo := generateEcho("get_forward_msg")
	data := map[string]interface{}{
		"action": "get_forward_msg",
		"params": map[string]interface{}{
			"id": id,
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
		return nil, fmt.Errorf("获取合并转发消息失败: %s", resp.Status)
	}

	var forwardMsg ForwardMessage
	err = json.Unmarshal(resp.Data, &forwardMsg)
	if err != nil {
		return nil, err
	}

	return &forwardMsg, nil
}

// SendLike 发送好友赞
func SendLike(b plugin.Bot, userID int64, times int) error {
	if times <= 0 {
		times = 1
	}
	if times > 10 {
		times = 10
	}

	echo := generateEcho("send_like")
	data := map[string]interface{}{
		"action": "send_like",
		"params": map[string]interface{}{
			"user_id": userID,
			"times":   times,
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
		return fmt.Errorf("点赞失败: %s", resp.Status)
	}

	return nil
}

// GetStrangerInfo 获取陌生人信息
func GetStrangerInfo(b plugin.Bot, userID int64) (*StrangerInfo, error) {
	// 生成唯一 echo
	echo := generateEcho("get_stranger_info")

	// 构造请求数据
	data := map[string]interface{}{
		"action": "get_stranger_info",
		"params": map[string]interface{}{
			"user_id": userID,
		},
		"echo": echo,
	}

	// 发送请求
	err := b.WriteJSON(data)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	// 等待响应
	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return nil, err
	}

	// 检查响应状态
	if resp.Status != "ok" {
		return nil, fmt.Errorf("API 调用失败: %s (retcode: %d)", resp.Status, resp.Retcode)
	}

	// 解析用户信息
	var strangerInfo StrangerInfo
	err = json.Unmarshal(resp.Data, &strangerInfo)
	if err != nil {
		return nil, fmt.Errorf("解析响应数据失败: %v", err)
	}

	return &strangerInfo, nil
}

// GetGroupInfo 获取群信息
func GetGroupInfo(b plugin.Bot, groupID int64) (*GroupInfo, error) {
	echo := generateEcho("get_group_info")

	data := map[string]interface{}{
		"action": "get_group_info",
		"params": map[string]interface{}{
			"group_id": groupID,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return nil, err
	}

	if resp.Status != "ok" {
		return nil, fmt.Errorf("API 调用失败: %s (retcode: %d)", resp.Status, resp.Retcode)
	}

	var groupInfo GroupInfo
	err = json.Unmarshal(resp.Data, &groupInfo)
	if err != nil {
		return nil, fmt.Errorf("解析响应数据失败: %v", err)
	}

	return &groupInfo, nil
}
