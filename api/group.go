package api

import (
	"encoding/json"
	"fmt"

	"github.com/iamlibie/milonra-go/plugin"
)

// SetGroupKick 群组踢人
func SetGroupKick(b plugin.Bot, groupID, userID int64, rejectAddRequest bool) error {
	echo := generateEcho("set_group_kick")
	data := map[string]interface{}{
		"action": "set_group_kick",
		"params": map[string]interface{}{
			"group_id":           groupID,
			"user_id":            userID,
			"reject_add_request": rejectAddRequest,
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
		return fmt.Errorf("踢人失败: %s", resp.Status)
	}

	return nil
}

// SetGroupBan 群组单人禁言
func SetGroupBan(b plugin.Bot, groupID, userID int64, duration int) error {
	echo := generateEcho("set_group_ban")
	data := map[string]interface{}{
		"action": "set_group_ban",
		"params": map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
			"duration": duration,
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
		return fmt.Errorf("禁言失败: %s", resp.Status)
	}

	return nil
}

// SetGroupAnonymousBan 群组匿名用户禁言
func SetGroupAnonymousBan(b plugin.Bot, groupID string, flag string, duration int) error {
	echo := generateEcho("set_group_anonymous_ban")
	data := map[string]interface{}{
		"action": "set_group_anonymous_ban",
		"params": map[string]interface{}{
			"group_id": groupID,
			"flag":     flag,
			"duration": duration,
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
		return fmt.Errorf("匿名禁言失败: %s", resp.Status)
	}

	return nil
}

// SetGroupWholeBan 群组全员禁言
func SetGroupWholeBan(b plugin.Bot, groupID int64, enable bool) error {
	echo := generateEcho("set_group_whole_ban")
	data := map[string]interface{}{
		"action": "set_group_whole_ban",
		"params": map[string]interface{}{
			"group_id": groupID,
			"enable":   enable,
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
		return fmt.Errorf("全员禁言失败: %s", resp.Status)
	}

	return nil
}

// SetGroupAdmin 设置群管理员
func SetGroupAdmin(b plugin.Bot, groupID, userID int64, enable bool) error {
	echo := generateEcho("set_group_admin")
	data := map[string]interface{}{
		"action": "set_group_admin",
		"params": map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
			"enable":   enable,
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
		return fmt.Errorf("设置管理员失败: %s", resp.Status)
	}

	return nil
}

// SetGroupAnonymous 群组匿名
func SetGroupAnonymous(b plugin.Bot, groupID int64, enable bool) error {
	echo := generateEcho("set_group_anonymous")
	data := map[string]interface{}{
		"action": "set_group_anonymous",
		"params": map[string]interface{}{
			"group_id": groupID,
			"enable":   enable,
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
		return fmt.Errorf("设置匿名失败: %s", resp.Status)
	}

	return nil
}

// SetGroupCard 设置群名片
func SetGroupCard(b plugin.Bot, groupID, userID int64, card string) error {
	echo := generateEcho("set_group_card")
	data := map[string]interface{}{
		"action": "set_group_card",
		"params": map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
			"card":     card,
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
		return fmt.Errorf("设置群名片失败: %s", resp.Status)
	}

	return nil
}

// SetGroupName 设置群名
func SetGroupName(b plugin.Bot, groupID int64, name string) error {
	echo := generateEcho("set_group_name")
	data := map[string]interface{}{
		"action": "set_group_name",
		"params": map[string]interface{}{
			"group_id":   groupID,
			"group_name": name,
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
		return fmt.Errorf("设置群名失败: %s", resp.Status)
	}

	return nil
}

// SetGroupLeave 退群
func SetGroupLeave(b plugin.Bot, groupID int64, isDismiss bool) error {
	echo := generateEcho("set_group_leave")
	data := map[string]interface{}{
		"action": "set_group_leave",
		"params": map[string]interface{}{
			"group_id":   groupID,
			"is_dismiss": isDismiss,
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
		return fmt.Errorf("退群失败: %s", resp.Status)
	}

	return nil
}

// SetGroupSpecialTitle 设置群专属头衔
func SetGroupSpecialTitle(b plugin.Bot, groupID, userID int64, specialTitle string, duration int) error {
	echo := generateEcho("set_group_special_title")
	data := map[string]interface{}{
		"action": "set_group_special_title",
		"params": map[string]interface{}{
			"group_id":      groupID,
			"user_id":       userID,
			"special_title": specialTitle,
			"duration":      duration,
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
		return fmt.Errorf("设置专属头衔失败: %s", resp.Status)
	}

	return nil
}

// GetGroupMemberInfo 获取群成员信息
func GetGroupMemberInfo(b plugin.Bot, groupID, userID int64, noCache bool) (*GroupMemberInfo, error) {
	echo := generateEcho("get_group_member_info")
	data := map[string]interface{}{
		"action": "get_group_member_info",
		"params": map[string]interface{}{
			"group_id": groupID,
			"user_id":  userID,
			"no_cache": noCache,
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
		return nil, fmt.Errorf("获取群成员信息失败: %s", resp.Status)
	}

	var memberInfo GroupMemberInfo
	err = json.Unmarshal(resp.Data, &memberInfo)
	if err != nil {
		return nil, err
	}

	return &memberInfo, nil
}
