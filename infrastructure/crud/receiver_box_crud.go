package crud

import (
	"time"
	"tiny_talk/infrastructure/db"
	"tiny_talk/infrastructure/models"
)

type recvBoxCRUD struct{}

var RecvBoxCRUD *recvBoxCRUD

func (*recvBoxCRUD) Create(message *models.ReceiverBox) error {
	return db.MysqlClient.Create(message).Error
}

func (*recvBoxCRUD) Update(message *models.ReceiverBox) error {
	return db.MysqlClient.Where("user_id = ? AND msg_id = ?", message.UserId, message.MsgId).Updates(message).Error
}

func (*recvBoxCRUD) GetmessageList(userId int64, time time.Time) ([]*models.MessageBasic, error) {
	var messages []*models.MessageBasic
	err := db.MysqlClient.
		Select("message_basic.*").
		Table("receiver_box").
		Joins("INNER JOIN message_basic ON receiver_box.msg_id = message_basic.id").
		Where("receiver_box.user_id = ? AND receiver_box.send_time > ?", userId, time).
		Order("receiver_box.send_time ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (*recvBoxCRUD) Delete(message *models.ReceiverBox) error {
	return db.MysqlClient.Where("user_id = ? AND msg_id = ?", message.UserId, message.MsgId).Delete(message).Error
}
