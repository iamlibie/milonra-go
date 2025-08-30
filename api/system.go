package api

import (
	"encoding/json"
	"fmt"

	"github.com/iamlibie/milonra-go/plugin"
)

// SetFriendAddRequest 处理加好友请求
func SetFriendAddRequest(b plugin.Bot, flag string, approve bool, remark string) error {
	echo := generateEcho("set_friend_add_request")
	data := map[string]interface{}{
		"action": "set_friend_add_request",
		"params": map[string]interface{}{
			"flag":    flag,
			"approve": approve,
			"remark":  remark,
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
		return fmt.Errorf("处理好友请求失败: %s", resp.Status)
	}

	return nil
}

// SetGroupAddRequest 处理加群请求/邀请
func SetGroupAddRequest(b plugin.Bot, flag, subType string, approve bool, reason string) error {
	echo := generateEcho("set_group_add_request")
	data := map[string]interface{}{
		"action": "set_group_add_request",
		"params": map[string]interface{}{
			"flag":     flag,
			"sub_type": subType,
			"approve":  approve,
			"reason":   reason,
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
		return fmt.Errorf("处理群请求失败: %s", resp.Status)
	}

	return nil
}

// GetLoginInfo 获取登录号信息
func GetLoginInfo(b plugin.Bot) (*LoginInfo, error) {
	echo := generateEcho("get_login_info")
	data := map[string]interface{}{
		"action": "get_login_info",
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
		return nil, fmt.Errorf("获取登录信息失败: %s", resp.Status)
	}

	var loginInfo LoginInfo
	err = json.Unmarshal(resp.Data, &loginInfo)
	if err != nil {
		return nil, err
	}

	return &loginInfo, nil
}

// GetFriendList 获取好友列表
func GetFriendList(b plugin.Bot) ([]FriendInfo, error) {
	echo := generateEcho("get_friend_list")
	data := map[string]interface{}{
		"action": "get_friend_list",
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
		return nil, fmt.Errorf("获取好友列表失败: %s", resp.Status)
	}

	var friends []FriendInfo
	err = json.Unmarshal(resp.Data, &friends)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

// GetCookies 获取Cookies
func GetCookies(b plugin.Bot, domain string) (*Cookies, error) {
	echo := generateEcho("get_cookies")
	data := map[string]interface{}{
		"action": "get_cookies",
		"params": map[string]interface{}{
			"domain": domain,
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
		return nil, fmt.Errorf("获取Cookies失败: %s", resp.Status)
	}

	var cookies Cookies
	err = json.Unmarshal(resp.Data, &cookies)
	if err != nil {
		return nil, err
	}

	return &cookies, nil
}

// GetCSRFToken 获取CSRF Token
func GetCSRFToken(b plugin.Bot) (*CSRFToken, error) {
	echo := generateEcho("get_csrf_token")
	data := map[string]interface{}{
		"action": "get_csrf_token",
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
		return nil, fmt.Errorf("获取CSRF Token失败: %s", resp.Status)
	}

	var token CSRFToken
	err = json.Unmarshal(resp.Data, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetCredentials 获取QQ相关接口凭证
func GetCredentials(b plugin.Bot, domain string) (*Credentials, error) {
	echo := generateEcho("get_credentials")
	data := map[string]interface{}{
		"action": "get_credentials",
		"params": map[string]interface{}{
			"domain": domain,
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
		return nil, fmt.Errorf("获取凭证失败: %s", resp.Status)
	}

	var credentials Credentials
	err = json.Unmarshal(resp.Data, &credentials)
	if err != nil {
		return nil, err
	}

	return &credentials, nil
}

// GetRecord 获取语音
func GetRecord(b plugin.Bot, file, outFormat string) (*RecordInfo, error) {
	echo := generateEcho("get_record")
	data := map[string]interface{}{
		"action": "get_record",
		"params": map[string]interface{}{
			"file":       file,
			"out_format": outFormat,
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
		return nil, fmt.Errorf("获取语音失败: %s", resp.Status)
	}

	var recordInfo RecordInfo
	err = json.Unmarshal(resp.Data, &recordInfo)
	if err != nil {
		return nil, err
	}

	return &recordInfo, nil
}

// GetImage 获取图片
func GetImage(b plugin.Bot, file string) (*ImageInfo, error) {
	echo := generateEcho("get_image")
	data := map[string]interface{}{
		"action": "get_image",
		"params": map[string]interface{}{
			"file": file,
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
		return nil, fmt.Errorf("获取图片失败: %s", resp.Status)
	}

	var imageInfo ImageInfo
	err = json.Unmarshal(resp.Data, &imageInfo)
	if err != nil {
		return nil, err
	}

	return &imageInfo, nil
}

// CanSendImage 检查是否可以发送图片
func CanSendImage(b plugin.Bot) (bool, error) {
	echo := generateEcho("can_send_image")
	data := map[string]interface{}{
		"action": "can_send_image",
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return false, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return false, err
	}

	if resp.Status != "ok" {
		return false, fmt.Errorf("检查失败: %s", resp.Status)
	}

	var canSend CanSend
	err = json.Unmarshal(resp.Data, &canSend)
	if err != nil {
		return false, err
	}

	return canSend.Yes, nil
}

// CanSendRecord 检查是否可以发送语音
func CanSendRecord(b plugin.Bot) (bool, error) {
	echo := generateEcho("can_send_record")
	data := map[string]interface{}{
		"action": "can_send_record",
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return false, err
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return false, err
	}

	if resp.Status != "ok" {
		return false, fmt.Errorf("检查失败: %s", resp.Status)
	}

	var canSend CanSend
	err = json.Unmarshal(resp.Data, &canSend)
	if err != nil {
		return false, err
	}

	return canSend.Yes, nil
}

// GetStatus 获取运行状态
func GetStatus(b plugin.Bot) (*Status, error) {
	echo := generateEcho("get_status")
	data := map[string]interface{}{
		"action": "get_status",
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
		return nil, fmt.Errorf("获取状态失败: %s", resp.Status)
	}

	var status Status
	err = json.Unmarshal(resp.Data, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// GetVersionInfo 获取版本信息
func GetVersionInfo(b plugin.Bot) (*VersionInfo, error) {
	echo := generateEcho("get_version_info")
	data := map[string]interface{}{
		"action": "get_version_info",
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
		return nil, fmt.Errorf("获取版本信息失败: %s", resp.Status)
	}

	var versionInfo VersionInfo
	err = json.Unmarshal(resp.Data, &versionInfo)
	if err != nil {
		return nil, err
	}

	return &versionInfo, nil
}

// SetRestart 重启OneBot实现
func SetRestart(b plugin.Bot, delay int) error {
	echo := generateEcho("set_restart")
	data := map[string]interface{}{
		"action": "set_restart",
		"params": map[string]interface{}{
			"delay": delay,
		},
		"echo": echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		return err
	}

	// 重启是异步操作，可能不会收到响应
	// 设置较短的超时时间
	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		// 如果是超时错误，可能是正常的（因为服务正在重启）
		return nil
	}

	if resp.Status != "ok" && resp.Status != "async" {
		return fmt.Errorf("重启失败: %s", resp.Status)
	}

	return nil
}

// CleanCache 清理缓存
func CleanCache(b plugin.Bot) error {
	echo := generateEcho("clean_cache")
	data := map[string]interface{}{
		"action": "clean_cache",
		"echo":   echo,
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
		return fmt.Errorf("清理缓存失败: %s", resp.Status)
	}

	return nil
}
