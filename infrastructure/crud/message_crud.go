package crud

import (
	"time"
	"tiny_talk/infrastructure/db"
	"tiny_talk/infrastructure/models"
)

type messageCRUD struct{}

var MessageCRUD *messageCRUD

func (*messageCRUD) Create(message *models.MessageBasic) error {
	return db.MysqlClient.Create(message).Error
}

func (*messageCRUD) Update(message *models.MessageBasic) error {
	return db.MysqlClient.Where("user_id = ? AND dest_id = ?", message.UserId, message.DestId).Updates(message).Error
}

func (*messageCRUD) Get(userId int64, seqNum int64) (*models.MessageBasic, error) {
	var message *models.MessageBasic
	err := db.MysqlClient.Where("user_id = ? AND seq_num = ?", userId, seqNum).Find(&message).Error
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (*messageCRUD) GetmessageList(userId int64, destId int64, time time.Time) ([]*models.MessageBasic, error) {
	var messages []*models.MessageBasic
	err := db.MysqlClient.Where("user_id = ? AND dest_id = ? AND send_time > ?", userId, destId, time).
		Order("send_time ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (*messageCRUD) Delete(message *models.MessageBasic) error {
	return db.MysqlClient.Where("user_id = ? AND seq_num = ?", message.UserId, message.SeqNum).Delete(message).Error
}
