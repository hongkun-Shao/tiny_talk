package utils

import (
	"fmt"
	"tiny_talk/utils/config"

	"github.com/bwmarrin/snowflake"
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
