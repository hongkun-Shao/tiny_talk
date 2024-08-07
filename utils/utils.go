package utils

import (
	"encoding/base64"
	"fmt"
	"tiny_talk/utils/config"

	"github.com/bwmarrin/snowflake"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

func GetConfigPath() string {
	path := "utils/config/config.toml"
	return path
}

func ParseToDsn(cfg *config.MysqlConfig) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.DbAddress, cfg.DbName)
	return dsn
}

func GenerateSnowflakeID() int64 {
	node, _ := snowflake.NewNode(1)
	return node.Generate().Int64()
}

func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword hashes the password with a cost of 10.
	// The cost parameter can be adjusted for more security but will also increase computation time.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword string, candidatePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatePassword))
	return err == nil
}

func GenerateToken() (string, error) {
	b := make([]byte, 32) //32 byte rand data
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
