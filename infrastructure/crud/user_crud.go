package crud

import (
	"tiny_talk/infrastructure/db"
	"tiny_talk/infrastructure/models"
)

type userCRUD struct{}

var UserCRUD *userCRUD

func (*userCRUD) Create(user *models.UserBasic) error {
	return db.MysqlClient.Create(user).Error
}

func (*userCRUD) Update(user *models.UserBasic) error {
	return db.MysqlClient.Where("identity = ?", user.Identity).Updates(user).Error
}

func (*userCRUD) Get(userId int64) (*models.UserBasic, error) {
	var user = &models.UserBasic{
		Identity: userId,
	}
	if err := db.MysqlClient.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (*userCRUD) GetByIds(userId []int64) ([]*models.UserBasic, error) {
	var users []*models.UserBasic
	err := db.MysqlClient.Where("identity IN (?)", userId).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (*userCRUD) Delete(user *models.UserBasic) error {
	return db.MysqlClient.Where("identity = ? AND deleted_at IS NULL", user.Identity).Delete(user).Error
}
