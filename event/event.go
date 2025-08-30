package event

type MessageEvent struct {
	GroupID    int64                  `json:"group_id"`
	UserID     int64                  `json:"user_id"`
	Message    string                 `json:"message"`
	RawMessage string                 `json:"raw_message"` // 原始消息（带 CQ 码）
	Nickname   string                 `json:"nickname"`
	Time       int64                  `json:"time"`
	IsAtMe     bool                   `json:"is_at_me"` // 是否 @ 了机器人
	RawData    map[string]interface{} // 原始 JSON 数据
}
