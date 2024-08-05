package service

import "github.com/gin-gonic/gin"

// @Description index
// @Tags Index
// @Accept json
// @Produce json
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "welcome",
	})
}
