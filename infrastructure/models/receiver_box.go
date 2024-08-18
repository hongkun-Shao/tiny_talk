package models

import (
	"time"

	"gorm.io/gorm"
)

type ReceiverBox struct {
	gorm.Model
	UserId   int64      `gorm:"type:bigint; not null;"`
	MsgId    uint       `gorm:"type:bigint; not null;"`
	SendTime *time.Time `gorm:"type:datetime; not null;"`
}

func (table *ReceiverBox) TableName() string {
	return "receiver_box"
}
