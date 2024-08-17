package models

import (
	"time"

	"gorm.io/gorm"
)

type MessageBasic struct {
	gorm.Model
	UserId   int64      `gorm:"type:bigint; not null;"`
	DestId   int64      `gorm:"type:bigint; not null;"`
	Type     int8       `gorm:"type:tinyint(4);not null"` // 0. 未定义 1. 单聊消息 2. 群聊消息 3. 好友请求 4.广播消息
	Status   int8       `gorm:"type:tinyint(4);not null"` // 0. 未处理 1. 已发送 2. 已接收 3. 已读
	Content  string     `gorm:"type:varchar(255);"`
	SeqNum   int64      `gorm:"type:bigint; not null;"`
	SendTime *time.Time `gorm:"type:datetime;"`
}

func (table *MessageBasic) TableName() string {
	return "message_basic"
}
