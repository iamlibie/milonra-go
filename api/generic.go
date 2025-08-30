package api

import (
	"encoding/json"
	"fmt"

	"github.com/iamlibie/milonra-go/plugin"
)

// OCRText OCR识别结果文本
type OCRText struct {
	Text        string     `json:"text"`        // 识别的文本
	Confidence  int        `json:"confidence"`  // 置信度
	Coordinates []OCRPoint `json:"coordinates"` // 文本位置坐标
}

// OCRPoint OCR坐标点
type OCRPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// OCRResult OCR识别结果
type OCRResult struct {
	Texts    []OCRText `json:"texts"`    // 识别的文本列表
	Language string    `json:"language"` // 识别的语言
}

// OCRImage 图像OCR识别
func OCRImage(b plugin.Bot, image string) (*OCRResult, error) {
	echo := generateEcho("ocr_image")
	params := map[string]interface{}{
		"image": image,
	}
	data := map[string]interface{}{
		"action": "ocr_image",
		"params": params,
		"echo":   echo,
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
		return nil, fmt.Errorf("OCR识别失败: %s", resp.Status)
	}

	var ocrResult OCRResult
	if err := json.Unmarshal(resp.Data, &ocrResult); err != nil {
		return nil, fmt.Errorf("解析OCR结果失败: %v", err)
	}

	return &ocrResult, nil
}

// SendGroupPoke 发送群聊戳一戳
func SendGroupPoke(b plugin.Bot, groupID int64, userID int64) error {
	echo := generateEcho("group_poke")
	params := map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	}
	data := map[string]interface{}{
		"action": "group_poke",
		"params": params,
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		// 如果失败，尝试使用消息方式
		msg := NewMessage().Poke(userID)
		_, err2 := SendGroupMessage(b, groupID, msg)
		return err2
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("戳一戳失败: %s", resp.Status)
	}

	return nil
}

// SendPrivatePoke 发送私聊戳一戳
func SendPrivatePoke(b plugin.Bot, userID int64, targetID int64) error {
	echo := generateEcho("friend_poke")
	params := map[string]interface{}{
		"user_id": targetID,
	}
	data := map[string]interface{}{
		"action": "friend_poke",
		"params": params,
		"echo":   echo,
	}

	err := b.WriteJSON(data)
	if err != nil {
		// 如果失败，尝试使用消息方式
		msg := NewMessage().Poke(targetID)
		_, err2 := SendPrivateMessage(b, targetID, msg)
		return err2
	}

	resp, err := responseWaiter.WaitForResponse(echo)
	if err != nil {
		return err
	}

	if resp.Status != "ok" {
		return fmt.Errorf("戳一戳失败: %s", resp.Status)
	}

	return nil
}
