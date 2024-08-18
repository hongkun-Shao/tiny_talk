package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"tiny_talk/infrastructure/crud"
	"tiny_talk/infrastructure/db"
	"tiny_talk/infrastructure/models"
	"tiny_talk/utils/logger"
)

type Message struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Type       int8   `json:"type"` // 0. 未定义 1. 单聊消息 2. 群聊消息 3. 好友请求 4.广播消息
	// Status   int8       `gorm:"type:tinyint(4);not null"` // 0. 未处理 1. 已发送 2. 已接收 3. 已读
	Content string `json:"content"`
}

func Publish(uid int64, message string) error {
	ctx := context.Background()
	channel := fmt.Sprintf("recv_box:%d", uid)
	err := db.Push(&ctx, channel, message).Err()
	if err != nil {
		logger.Infof("Publish err: %v", err)
		return err
	}
	err = db.Expire(&ctx, channel, 7*24*time.Hour).Err()
	if err != nil {
		logger.Infof("Expire err: %v", err)
		return err
	}
	logger.Info("Publish success")
	return nil
}

func Subcribe(uid int64) error {
	ctx := context.Background()
	channel := fmt.Sprintf("recv_box:%d", uid)
	OnlineUsers.Locker.RLock()
	targetUser := OnlineUsers.Connections[uid]
	OnlineUsers.Locker.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	length, err := db.Len(&ctx, channel).Result()
	if err != nil {
		logger.Fatalf("Failed to get message from queue: %v", err)
		return err
	}

	// mysql offline message
	if length == 0 {
		logger.Infof("No messages in the queue, pull messages from mysql")
		var messages []*models.MessageBasic
		start_time := time.Now().Add(-7 * 24 * time.Hour)
		messages, _ = crud.RecvBoxCRUD.GetmessageList(uid, start_time)
		for _, message_basic := range messages {
			message := Message{
				SenderID:   string(message_basic.UserId),
				ReceiverID: string(message_basic.DestId),
				Type:       message_basic.Type,
				Content:    message_basic.Content,
			}
			messageJSON, _ := json.Marshal(message)
			targetUser.Channel <- string(messageJSON)
		}
		return nil
	}

	//redis offline message
	for {
		msg, err := db.Pop(&ctx, 0, channel).Result()
		if err != nil {
			if err == context.DeadlineExceeded {
				logger.Infof("No more messages in the queue.")
				break
			}
			logger.Fatalf("Failed to get message from queue: %v", err)
			return err
		}
		logger.Infof("Received message '%s' from the queue.\n", msg[1])
		targetUser.Channel <- msg[1]
		logger.Infof("Send message to client")
	}
	return nil
}
