package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"tiny_talk/infrastructure/crud"
	"tiny_talk/infrastructure/db"
	"tiny_talk/infrastructure/models"
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
	user.Name = c.PostForm("username")
	user.Password = c.PostForm("password")
	repassword := c.PostForm("repassword")
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

	err = crud.UserCRUD.Create(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		logger.Errorf("Failed to create user: %v", err)
		return
	}

	c.JSON(200, gin.H{
		"message":  "User created successfully, please remember your identity",
		"username": user.Name,
		"identity": fmt.Sprintf("%d", user.Identity),
	})
	logger.Infof("User created successfully: %v", user.Identity)
}

// @Description login
// @Tags User
// @Param identity formData string true "账号"
// @Param password formData string true "密码"
// @Accept json
// @Produce json
// @Success 200 {string} create user
// @Router /user/Login [post]
func UserLogin(c *gin.Context) {
	var user *models.UserBasic
	identity, err := strconv.ParseInt(c.PostForm("userid"), 10, 64)
	password := c.PostForm("password")
	if err != nil {
		c.JSON(400, gin.H{"error": "Identity must be a number"})
		logger.Errorf("Identity must be a number: %v", err)
		return
	}
	user, err = crud.UserCRUD.Get(identity)

	res := utils.CheckPassword(user.Password, password)
	if !res {
		c.JSON(400, gin.H{"error": "Password is incorrect"})
		logger.Errorf("Password is incorrect")
		return
	}
	token, err := utils.GenerateToken()
	// 成功才存储 token, 否则返回空的token
	if err == nil {
		err = storeTokenInRedis(token, user.Identity)
		// 存储失败，就返回空的 token
		if err != nil {
			logger.Errorf("Login store token failed %v", err)
			token = ""
		}
	} else {
		logger.Errorf("Login generate token failed %v", err)
	}
	// TODO : store user in redis

	c.JSON(200, gin.H{
		"message": "Login successfully",
		"token":   token,
		// "redirect": "webui/dashboard",
	})
	logger.Infof("Login successfully: id %v, token %v", user.Identity, token)
}

// @Description test token
// @Tags User
// @Param token query string false "密令"
// @Accept json
// @Produce json
// @Success 200 {string} test token
// @Router /user/TestToken [post]
func TestToken(c *gin.Context) {
	token := c.Query("token")
	id, err := GetUserIdFromRedisByToken(token)
	if err != nil {
		c.JSON(400, gin.H{
			"error":   "Token is invalid",
			"message": "Please login again",
		})
		logger.Errorf("Token is invalid not fount in redis")
		return
	}
	logger.Infof("TestToken: id %v, token %v", id, token)
	token, err = RefeshTokenExpiration(token)
	if err != nil {
		logger.Errorf("Token refresh expiration failed %v", err)
	}
	c.JSON(200, gin.H{
		"message": "Token is valid",
		"token":   token,
	})
}

func storeTokenInRedis(token string, identity int64) error {
	// 设置 token 的过期时间, 存活时间 3 hours, 但是最长存活时间 24 hours
	ctx := context.Background()
	err := db.Set(&ctx, token, identity, 3*time.Hour).Err()
	if err != nil {
		return errors.New("failed to store token in Redis")
	}
	now := time.Now().UnixMilli()
	// 存储 token 的创建时间, 存活24 + 3 + 1 hours, 用来检测 token 是否 超过 最长存活时间
	err = db.Set(&ctx, fmt.Sprintf("created:%s", token), now, 3*time.Hour).Err()
	if err != nil {
		return errors.New("failed to store token created time in Redis")
	}
	return nil
}

func deleteTokenFromRedis(token string) error {
	ctx := context.Background()
	err := db.Del(&ctx, token).Err()
	if err != nil {
		return errors.New("failed to delete token in Redis")
	}
	err = db.Del(&ctx, fmt.Sprintf("created:%s", token)).Err()
	if err != nil {
		return errors.New("failed to delete token created time in Redis")
	}
	return nil
}

func GetUserIdFromRedisByToken(token string) (int64, error) {
	ctx := context.Background()
	val, err := db.Get(&ctx, token).Result()
	if err != nil {
		return 0, err
	}
	id, _ := strconv.ParseInt(val, 10, 64)
	return id, nil
}

func getTokenCreatedTImeFromRedisByToken(token string) (int64, error) {
	ctx := context.Background()
	val, err := db.Get(&ctx, fmt.Sprintf("created:%s", token)).Result()
	if err != nil {
		return 0, err
	}
	created_time, _ := strconv.ParseInt(val, 10, 64)
	return created_time, nil
}

func RefeshTokenExpiration(token string) (string, error) {
	ctx := context.Background()
	err := db.Expire(&ctx, token, 3*time.Hour).Err()
	if err != nil {
		return token, err
	}
	now := time.Now().UnixMilli()

	// 不关心 err, 调用失败，会返回created_time = 0
	created_time, _ := getTokenCreatedTImeFromRedisByToken(token)
	if now-created_time >= 24*3600*1000 {
		logger.Infof("token expired, refresh token: %v", token)
		// TODO: 增加失败重试机制，避免 new_token 为空
		new_token, _ := utils.GenerateToken()
		logger.Infof("new token: %v", new_token)

		// TODO: 增加失败重试机制，避免 identeity 为空
		identity, _ := GetUserIdFromRedisByToken(token)
		logger.Infof("new token's identity: %v", identity)

		// TODO: 增加失败重试机制，避免删除 old_token 失败
		_ = deleteTokenFromRedis(token)

		// 设置失败，就返回空字符串，下次用户操作需要重新登录
		err = storeTokenInRedis(new_token, identity)
		if err != nil {
			return "", errors.New("failed to store new token in Redis")
		}
		return new_token, nil
	}
	return token, nil
}
