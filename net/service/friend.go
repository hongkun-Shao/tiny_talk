package service

import (
	"strconv"
	"tiny_talk/infrastructure/mysql"
	"tiny_talk/infrastructure/mysql/models"
	"tiny_talk/utils/logger"

	"github.com/gin-gonic/gin"
)

// @Description send application(make frind)
// @Tags Friend
// @Param friend_id query string true "对方账号"
// @Param token query string true "密令"
// @Accept json
// @Produce json
// @Success 200 {string} send application success
// @Router /friend/MakeFriendById [post]
func MakeFriendById(c *gin.Context) {
	var friend models.UserBasic
	var contact models.FriendBasic
	friend_id, err := strconv.ParseInt(c.Query("friend_id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Invalid friend_id",
			"message": "Please provide a valid friend_id",
		})
	}
	token := c.Query("token")

	user_id, err := GetUserIdFromRedisByToken(token)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Token is invalid",
			"message": "Please login again",
		})
		logger.Errorf("Token is invalid not fount in redis")
		return
	}

	if friend_id == user_id {
		c.JSON(400, gin.H{
			"error":   "Invalid friend_id",
			"message": "You can't make friend with yourself",
		})
		return
	}

	friend.Identity = friend_id
	err = mysql.FindModel(&friend)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "friend not exist",
			"message": "Please provide a valid friend_id",
		})
	}

	contact.UserId = user_id
	contact.FriendId = friend_id
	contact.Status = 1
	err = mysql.CreateModel(&contact)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to send application"})
		logger.Errorf("Failed to create friend contact: %v", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Application has been sent",
	})
	logger.Infof("user %v want to make friend with user %v", user_id, friend_id)
}
