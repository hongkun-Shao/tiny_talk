package service

import (
	"strconv"
	"tiny_talk/infrastructure/mysql"
	"tiny_talk/infrastructure/mysql/models"
	"tiny_talk/utils"
	"tiny_talk/utils/logger"

	"github.com/gin-gonic/gin"
)

// @Description create user
// @Tags User
// @Param name query string true "用户名"
// @Param password query string true "密码"
// @Param repassword query string true "确认密码"
// @Accept json
// @Produce json
// @Success 200 {string} create user
// @Router /user/CreateUser [post]
func CreateUser(c *gin.Context) {
	var user models.UserBasic
	user.Name = c.Query("name")
	user.Password = c.Query("password")
	repassword := c.Query("repassword")
	logger.Infof("%v, %v, %v", user.Name, user.Password, repassword)
	if user.Password != repassword {
		c.JSON(400, gin.H{"error": "Password and Repassword must be the same"})
		logger.Errorf("Password and Repassword must be the same")
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		logger.Errorf("Failed to hash password: %v", err)
		return
	}
	user.Password = hashedPassword
	user.Identity = utils.GenerateSnowflakeID()

	err = mysql.CreateModel(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		logger.Errorf("Failed to create user: %v", err)
		return
	}

	c.JSON(200, gin.H{
		"message":  "User created successfully, please remember your identity",
		"username": user.Name,
		"identity": user.Identity,
	})

	logger.Infof("User created successfully: %v", user.Identity)
}

// @Description login
// @Tags User
// @Param identity query string true "账号"
// @Param password query string true "密码"
// @Accept json
// @Produce json
// @Success 200 {string} create user
// @Router /user/Login [post]
func Login(c *gin.Context) {
	var user models.UserBasic
	identity, err := strconv.ParseInt(c.Query("identity"), 10, 64)
	password := c.Query("password")
	if err != nil {
		c.JSON(400, gin.H{"error": "Identity must be a number"})
		logger.Errorf("Identity must be a number: %v", err)
		return
	}
	user.Identity = identity
	mysql.FindModel(&user)
	res := utils.CheckPassword(user.Password, password)
	if !res {
		c.JSON(400, gin.H{"error": "Password is incorrect"})
		logger.Errorf("Password is incorrect")
		return
	}
	c.JSON(200, gin.H{
		"message": "Login successfully",
	})
	logger.Infof("Login successfully: %v", user.Identity)
}
