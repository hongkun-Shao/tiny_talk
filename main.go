package main

import (
	"log"
	"tiny_talk/utils"
	"tiny_talk/utils/config"
	"tiny_talk/utils/logger"
)

func main() {
	// ctx := context.Background()
	log.Printf("tiny_talk start ...")
	path := utils.GetConfigPath()
	appConfig, err := config.LoadAppStaticConfig(path)
	if err != nil {
		log.Fatalf("failed to load config file, err = %v", err)
	}

	// 日志初始化
	_, err = logger.InitLogger(appConfig.Loggers, "text", "main")
	if err != nil {
		log.Fatalf("failed to init text format logger, err = %v", err)
	}
	// logger.Info(123456789)
	logger.Infof("message is %v", 1234567898)
	// for {
	// 	logger.Info("test log Rolling...")
	// }
}
