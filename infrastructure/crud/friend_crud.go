package crud

import (
	"tiny_talk/infrastructure/db"
	"tiny_talk/infrastructure/models"
)

type friendCRUD struct{}

var FriendCRUD *friendCRUD

func (*friendCRUD) Create(friend *models.FriendBasic) error {
	return db.MysqlClient.Create(friend).Error
}

func (*friendCRUD) Update(friend *models.FriendBasic) error {
	return db.MysqlClient.Where("user_id = ? AND friend_id = ?", friend.UserId, friend.FriendId).Updates(friend).Error
}

func (*friendCRUD) Get(userId int64, friendId int64) (*models.FriendBasic, error) {
	var friend *models.FriendBasic
	err := db.MysqlClient.Where("user_id = ? AND friend_id = ?", userId, friendId).Find(&friend).Error
	if err != nil {
		return nil, err
	}
	return friend, nil
}

func (*friendCRUD) GetFriendList(userId int64) ([]*models.FriendBasic, error) {
	var friends []*models.FriendBasic
	err := db.MysqlClient.Where("user_id = ?", userId).Find(&friends).Error
	if err != nil {
		return nil, err
	}
	return friends, nil
}

func (*friendCRUD) Delete(friend *models.FriendBasic) error {
	return db.MysqlClient.Where("user_id = ? AND friend_id = ?", friend.UserId, friend.FriendId).Delete(friend).Error
}
