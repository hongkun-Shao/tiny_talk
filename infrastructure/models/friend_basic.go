package models

import (
	"gorm.io/gorm"
)

type FriendBasic struct {
	gorm.Model
	UserId   int64  `gorm:"type:bigint; not null; index:uk_user_id_friend_id"`
	FriendId int64  `gorm:"type:bigint; not null; index:uk_user_id_friend_id"`
	Remaks   string `gorm:"type:varchar(20);"`        // 备注
	Status   int8   `gorm:"type:tinyint(4);not null"` // 1. 申请中 2. 已同意 3. 已拒绝
}

func (table *FriendBasic) TableName() string {
	return "friend_basic"
}
