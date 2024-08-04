package config

import (
	"log"
	"tiny_talk/utils/logger"

	"github.com/BurntSushi/toml"
)

type MysqlConfig struct {
	DbName    string `toml:"dbName"`
	User      string `toml:"dbUser"`
	Password  string `toml:"dbPassWd"`
	DbAddress string `toml:"dbAddress"`
}

type AppConfig struct {
	AppName string                         `toml:"app_name"`
	Loggers map[string]logger.LoggerConfig `toml:"loggers"`
	Mysql   MysqlConfig                    `toml:"mysql"`
	// Redis RedisParam `toml:"redis"`
}

func LoadAppStaticConfig(tomlPath string) (*AppConfig, error) {
	var config AppConfig
	_, err := toml.DecodeFile(tomlPath, &config)
	if err != nil {
		log.Printf("decode file fail, tomlPath:%s, err:%v", tomlPath, err)
		return nil, err
	}
	return &config, err
}
