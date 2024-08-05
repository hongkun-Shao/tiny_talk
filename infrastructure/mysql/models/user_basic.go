package models

import (
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string     `gorm:"type:varchar(20);not null"`
	Password      string     `gorm:"type:varchar(60);not null"`
	Sex           int8       `gorm:"type:tinyint(4);not null"` // 0:未知, 1:男, 2:女
	Avatar        string     `gorm:"type:varchar(100);"`       // 头像url
	Phone         string     `gorm:"type:varchar(11);"`
	Email         string     `gorm:"type:varchar(50);"`
	Identity      int64      `gorm:"type:bigint;primarykey;not null"`
	ClientIp      string     `gorm:"type:varchar(20);"`
	ClientPort    string     `gorm:"type:varchar(20);"`
	LoginTime     *time.Time `gorm:"type:datetime;"`
	LogoutTime    *time.Time `gorm:"type:datetime;"`
	HeartbeatTime *time.Time `gorm:"type:datetime;"`
	IsLogOut      bool       `gorm:"type:tinyint(4);not null"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func (table *UserBasic) Create(db *gorm.DB) error {
	return db.Create(table).Error
}

func (table *UserBasic) Update(db *gorm.DB) error {
	return db.Where("identity = ?", table.Identity).Updates(table).Error
}

func (table *UserBasic) Find(db *gorm.DB) error {
	return db.Where("identity = ?", table.Identity).First(table).Error
}

func (table *UserBasic) Delete(db *gorm.DB) error {
	return db.Where("identity = ? AND deleted_at IS NULL", table.Identity).Delete(table).Error
}
