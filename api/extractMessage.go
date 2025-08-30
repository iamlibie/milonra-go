package api

import (
	"fmt"
	"strings"
)

// ExtractMessage 从 OneBot 消息中提取详细可读信息
func ExtractMessage(data map[string]interface{}) string {
	// 情况1：纯字符串消息（少见）
	if msg, ok := data["message"].(string); ok {
		return msg
	}

	// 情况2：消息段数组
	if segs, ok := data["message"].([]interface{}); ok {
		var parts []string
		for _, seg := range segs {
			if segMap, ok := seg.(map[string]interface{}); ok {
				parts = append(parts, parseSegmentDetailed(segMap))
			}
		}
		return JoinNonEmpty(parts, "\n")
	}

	return "[无法解析的消息]"
}

// parseSegmentDetailed 解析单个消息段，返回详细文本
func parseSegmentDetailed(seg map[string]interface{}) string {
	segType, ok := seg["type"].(string)
	if !ok {
		return "[未知消息]"
	}

	data, ok := seg["data"].(map[string]interface{})
	if !ok {
		return "[" + segType + "]"
	}

	switch segType {
	case "text":
		if text, ok := data["text"].(string); ok {
			return text
		}
		return ""

	case "at":
		if qq, ok := data["qq"].(string); ok {
			if qq == "all" {
				return "[@全体成员]"
			}
			return fmt.Sprintf("[@QQ:%s]", qq)
		}
		return "[@某人]"

	case "face":
		if id, ok := data["id"].(string); ok {
			return fmt.Sprintf("[QQ表情 ID:%s]", id)
		}
		return "[QQ表情]"

	case "mface":
		emojiID := GetString(data, "emoji_id")
		summary := GetString(data, "summary")
		if summary != "" {
			return fmt.Sprintf("[商城表情: %s]", summary)
		}
		if emojiID != "" {
			return fmt.Sprintf("[商城表情 ID:%s]", emojiID)
		}
		return "[商城表情]"

	case "image":
		file := GetString(data, "file")
		url := GetString(data, "url")
		summary := GetString(data, "summary")
		subType := GetInt(data, "sub_type")
		fileSize := GetInt64(data, "file_size")

		var info []string
		info = append(info, "[图片]")

		if summary != "" {
			info = append(info, "描述:"+summary)
		}
		if file != "" {
			info = append(info, "文件:"+file)
		}
		if url != "" {
			info = append(info, "URL:"+url)
		}
		if fileSize > 0 {
			sizeKB := fileSize / 1024
			if fileSize >= 1024*1024 {
				sizeKB = fileSize / 1024 / 1024
				info = append(info, fmt.Sprintf("大小:%dMB", sizeKB))
			} else {
				info = append(info, fmt.Sprintf("大小:%dKB", sizeKB))
			}
		}
		if subType > 0 {
			info = append(info, fmt.Sprintf("子类型:%d", subType))
		}

		return strings.Join(info, ", ")

	case "record":
		file := GetString(data, "file")
		url := GetString(data, "url")
		fileSize := GetInt64(data, "file_size")
		path := GetString(data, "path")

		var info []string
		info = append(info, "[语音]")

		if file != "" {
			info = append(info, "标识:"+file)
		}
		if url != "" {
			info = append(info, "URL:"+url)
		}
		if fileSize > 0 {
			sizeKB := fileSize / 1024
			info = append(info, fmt.Sprintf("大小:%dKB", sizeKB))
		}
		if path != "" {
			info = append(info, "路径:"+path)
		}

		return strings.Join(info, ", ")

	case "video":
		file := GetString(data, "file")
		url := GetString(data, "url")
		fileSize := GetInt64(data, "file_size")
		thumb := GetString(data, "thumb")

		var info []string
		info = append(info, "[视频]")

		if file != "" {
			info = append(info, "文件:"+file)
		}
		if url != "" {
			info = append(info, "在线:"+url)
		}
		if fileSize > 0 {
			sizeMB := fileSize / 1024 / 1024
			info = append(info, fmt.Sprintf("大小:%dMB", sizeMB))
		}
		if thumb != "" {
			info = append(info, "缩略图:"+thumb)
		}

		return strings.Join(info, ", ")

	case "file":
		name := GetString(data, "name")
		file := GetString(data, "file")
		fileID := GetString(data, "file_id")
		fileSize := GetInt64(data, "file_size")

		var info []string
		info = append(info, "[文件]")

		if name != "" {
			info = append(info, "名称:"+name)
		} else if file != "" {
			info = append(info, "文件:"+file)
		}
		if fileID != "" {
			info = append(info, "ID:"+fileID)
		}
		if fileSize > 0 {
			sizeKB := fileSize / 1024
			if fileSize >= 1024*1024 {
				sizeKB = fileSize / 1024 / 1024
				info = append(info, fmt.Sprintf("大小:%dMB", sizeKB))
			} else {
				info = append(info, fmt.Sprintf("大小:%dKB", sizeKB))
			}
		}

		return strings.Join(info, ", ")

	case "reply":
		if id, ok := data["id"].(string); ok {
			return fmt.Sprintf("[回复消息 ID:%s]", id)
		}
		return "[回复]"

	case "poke":
		id := GetString(data, "id")
		typ := GetString(data, "type")
		if id != "" {
			return fmt.Sprintf("[戳一戳 ID:%s,TYPE:%s]", id, typ)
		}
		return "[戳一戳]"

	case "dice":
		result := GetString(data, "result")
		return fmt.Sprintf("[骰子: %s]", result)

	case "rps":
		result := GetString(data, "result")
		return fmt.Sprintf("[猜拳: %s]", result)

	case "json":
		return "[卡片消息]"

	case "music":
		return "[音乐分享]"

	case "forward":
		return "[合并转发]"

	default:
		return fmt.Sprintf("[%s]", segType)
	}
}

func GetString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func GetInt(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return int64(f)
		}
	}
	return 0
}

func GetInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return int64(f)
		}
	}
	return 0
}

func JoinNonEmpty(parts []string, sep string) string {
	var nonEmpty []string
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	return strings.Join(nonEmpty, sep)
}
