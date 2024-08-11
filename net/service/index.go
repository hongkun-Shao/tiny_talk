package service

import "github.com/gin-gonic/gin"

// @Description index
// @Tags Index
// @Accept json
// @Produce json
// @Success 200 {string} welcome
// @Router / [get]
func GetIndex(c *gin.Context) {
	c.File("webui/index.html")
}

func Register(c *gin.Context) {
	c.File("webui/register.html")
}

func Home(c *gin.Context) {
	c.File("webui/home.html")
}
