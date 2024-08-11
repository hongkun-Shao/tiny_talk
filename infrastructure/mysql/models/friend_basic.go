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

func (table *FriendBasic) Create(db *gorm.DB) error {
	return db.Create(table).Error
}

func (table *FriendBasic) Update(db *gorm.DB) error {
	return db.Where("user_id = ? AND friend_id = ?", table.UserId, table.FriendId).Updates(table).Error
}

func (table *FriendBasic) Find(db *gorm.DB) error {
	return db.Where("user_id = ? AND friend_id = ?", table.UserId, table.FriendId).First(table).Error
}

func (table *FriendBasic) Delete(db *gorm.DB) error {
	return db.Where("user_id = ? AND friend_id = ? AND deleted_at IS NULL", table.UserId, table.FriendId).Delete(table).Error
}
