package service

import (
	"tiny_talk/infrastructure/crud"
	"tiny_talk/infrastructure/models"
	"tiny_talk/utils/logger"

	"github.com/gin-gonic/gin"
)

func makeFriendById(user_id int64, friend_id int64) bool {
	var contact models.FriendBasic

	if friend_id == user_id {
		logger.Infof("user %v want to make friend with himself", user_id)
		return false
	}

	_, err := crud.UserCRUD.Get(friend_id)

	if err != nil {
		logger.Errorf("Friend not exist%v", friend_id)
		return false
	}

	contact.UserId = user_id
	contact.FriendId = friend_id
	contact.Status = 1
	err = crud.FriendCRUD.Create(&contact)
	if err != nil {
		logger.Errorf("Failed to create friend contact: %v", err)
		return false
	}

	logger.Infof("user %v want to make friend with user %v", user_id, friend_id)
	return true
}

func GetFriendList(c *gin.Context) {
	token := c.GetHeader("Authorization")
	// 检查 Authorization 字段是否存在
	logger.Infof("GetFriendList token = %v", token)
	user_id, err := GetUserIdFromRedisByToken(token)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Token is invalid",
			"message": "Please login again",
		})
		logger.Errorf("Token is invalid not fount in redis")
		return
	}
	friends, err := crud.FriendCRUD.GetFriendList(user_id)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get friend list"})
		logger.Errorf("Failed to get friend list: %v", err)
		return
	}

	logger.Infof("GetFriendList user_id = %v", user_id)
	c.JSON(200, gin.H{
		"message": "Get friend list success",
		"friends": friends,
	})
}
