package service

type Message struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Type       int8   `json:"type"` // 0. 未定义 1. 单聊消息 2. 群聊消息 3. 好友请求 4.广播消息
	// Status   int8       `gorm:"type:tinyint(4);not null"` // 0. 未处理 1. 已发送 2. 已接收 3. 已读
	Content string `json:"content"`
}
