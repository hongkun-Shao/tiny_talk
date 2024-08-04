package config

import (
	"log"
	"tiny_talk/utils/logger"

	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	AppName string                         `toml:"app_name"`
	Loggers map[string]logger.LoggerConfig `toml:"loggers"`
	// Redis RedisParam
	// Mysql MysqlParam
}

func LoadAppStaticConfig(tomlPath string) (*AppConfig, error) {
	var config AppConfig
	_, err := toml.DecodeFile(tomlPath, &config)
	log.Print("config:", config)
	if err != nil {
		log.Printf("decode file fail, tomlPath:%s, err:%v", tomlPath, err)
		return nil, err
	}
	return &config, err
}
