package api

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// APIResponse OneBot API 响应结构
type APIResponse struct {
	Status  string          `json:"status"`
	Retcode int64           `json:"retcode"`
	Data    json.RawMessage `json:"data"`
	Echo    string          `json:"echo"`
}

// StrangerInfo 陌生人信息结构
type StrangerInfo struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex"`
	Age      int    `json:"age"`
}

// GroupInfo 群信息结构
type GroupInfo struct {
	GroupID        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	MemberCount    int64  `json:"member_count"`
	MaxMemberCount int64  `json:"max_member_count"`
}

// ResponseWaiter 响应等待器
type ResponseWaiter struct {
	mu      sync.RWMutex
	waiters map[string]chan *APIResponse
	timeout time.Duration
}

// 全局响应等待器
var responseWaiter = &ResponseWaiter{
	waiters: make(map[string]chan *APIResponse),
	timeout: 30 * time.Second, // 30秒超时
}

// WaitForResponse 等待指定 echo 的响应
func (rw *ResponseWaiter) WaitForResponse(echo string) (*APIResponse, error) {
	// 创建响应通道
	respChan := make(chan *APIResponse, 1)

	rw.mu.Lock()
	rw.waiters[echo] = respChan
	rw.mu.Unlock()

	// 确保清理
	defer func() {
		rw.mu.Lock()
		delete(rw.waiters, echo)
		rw.mu.Unlock()
		close(respChan)
	}()

	// 等待响应或超时
	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(rw.timeout):
		return nil, fmt.Errorf("等待响应超时: %s", echo)
	}
}

// HandleAPIResponse 处理 API 响应（应在消息处理循环中调用）
func (rw *ResponseWaiter) HandleAPIResponse(data map[string]interface{}) {
	// 检查是否是 API 响应
	if _, hasStatus := data["status"]; !hasStatus {
		return
	}

	// 解析响应
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	var resp APIResponse
	err = json.Unmarshal(jsonData, &resp)
	if err != nil {
		return
	}

	// 如果有 echo，发送给对应的等待者
	if resp.Echo != "" {
		rw.mu.RLock()
		if waiter, exists := rw.waiters[resp.Echo]; exists {
			select {
			case waiter <- &resp:
			default:
				// 通道已满，忽略
			}
		}
		rw.mu.RUnlock()
	}
}

// 生成唯一 echo
func generateEcho(action string) string {
	return fmt.Sprintf("%s_%d", action, time.Now().UnixNano())
}

// HandleAPIResponse 公共的 API 响应处理函数
func HandleAPIResponse(data map[string]interface{}) {
	responseWaiter.HandleAPIResponse(data)
}

// WaitForResponse 公共的响应等待函数
func WaitForResponse(echo string) (*APIResponse, error) {
	return responseWaiter.WaitForResponse(echo)
}
