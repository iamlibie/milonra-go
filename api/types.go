package api

import "encoding/json"

// Sender 发送者信息
type Sender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Sex      string `json:"sex,omitempty"`
	Age      int64  `json:"age,omitempty"`
	Card     string `json:"card,omitempty"`  // 群名片
	Area     string `json:"area,omitempty"`  // 地区
	Level    string `json:"level,omitempty"` // 成员等级
	Role     string `json:"role,omitempty"`  // 角色 owner/admin/member
	Title    string `json:"title,omitempty"` // 专属头衔
}

// Anonymous 匿名信息
type Anonymous struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Flag string `json:"flag"`
}

// File 文件信息
type File struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	Busid int64  `json:"busid"`
}

// Status 状态信息
type Status struct {
	Online bool `json:"online"`
	Good   bool `json:"good"`
}

// VersionInfo 版本信息
type VersionInfo struct {
	AppName         string `json:"app_name"`
	AppVersion      string `json:"app_version"`
	ProtocolVersion string `json:"protocol_version"`
}

// MessageInfo 消息信息
type MessageInfo struct {
	Time        int64           `json:"time"`
	MessageType string          `json:"message_type"`
	MessageID   int64           `json:"message_id"`
	RealID      int64           `json:"real_id"`
	Sender      Sender          `json:"sender"`
	Message     json.RawMessage `json:"message"`
}

// ForwardMessage 合并转发消息
type ForwardMessage struct {
	Message []MessageSegment `json:"message"`
}

// LoginInfo 登录信息
type LoginInfo struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
}

// FriendInfo 好友信息
type FriendInfo struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}

// GroupMemberInfo 群成员信息
type GroupMemberInfo struct {
	GroupID         int64  `json:"group_id"`
	UserID          int64  `json:"user_id"`
	Nickname        string `json:"nickname"`
	Card            string `json:"card"`
	Sex             string `json:"sex"`
	Age             int64  `json:"age"`
	Area            string `json:"area"`
	JoinTime        int64  `json:"join_time"`
	LastSentTime    int64  `json:"last_sent_time"`
	Level           string `json:"level"`
	Role            string `json:"role"`
	Unfriendly      bool   `json:"unfriendly"`
	Title           string `json:"title"`
	TitleExpireTime int64  `json:"title_expire_time"`
	CardChangeable  bool   `json:"card_changeable"`
}

// HonorInfo 群荣誉信息
type HonorInfo struct {
	GroupID          int64         `json:"group_id"`
	CurrentTalkative *HonorMember  `json:"current_talkative,omitempty"`
	TalkativeList    []HonorMember `json:"talkative_list,omitempty"`
	PerformerList    []HonorMember `json:"performer_list,omitempty"`
	LegendList       []HonorMember `json:"legend_list,omitempty"`
	StrongNewbieList []HonorMember `json:"strong_newbie_list,omitempty"`
	EmotionList      []HonorMember `json:"emotion_list,omitempty"`
}

// HonorMember 群荣誉成员
type HonorMember struct {
	UserID      int64  `json:"user_id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	DayCount    int64  `json:"day_count,omitempty"`
	Description string `json:"description,omitempty"`
}

// Cookies Cookie信息
type Cookies struct {
	Cookies string `json:"cookies"`
}

// CSRFToken CSRF令牌
type CSRFToken struct {
	Token int64 `json:"token"`
}

// Credentials 凭证信息
type Credentials struct {
	Cookies   string `json:"cookies"`
	CSRFToken int64  `json:"csrf_token"`
}

// RecordInfo 语音信息
type RecordInfo struct {
	File string `json:"file"`
}

// ImageInfo 图片信息
type ImageInfo struct {
	File string `json:"file"`
}

// CanSend 是否可发送
type CanSend struct {
	Yes bool `json:"yes"`
}

// APIError API错误
type APIError struct {
	Status  string `json:"status"`
	Retcode int64  `json:"retcode"`
	Msg     string `json:"msg,omitempty"`
	Wording string `json:"wording,omitempty"`
}

func (e *APIError) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	if e.Wording != "" {
		return e.Wording
	}
	return "API error"
}
