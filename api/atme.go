package api

import (
	"strconv"
	"strings"
)

func IsAtMe(data string, selfID int64) bool {
	// 检查 message 是否是数组
	if strings.Contains(data, "[@QQ:"+strconv.FormatInt(selfID, 10)+"]") {
		return true
	}
	return false
}
