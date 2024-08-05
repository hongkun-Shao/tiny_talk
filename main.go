package main

import (
	"log"
	"tiny_talk/infrastructure/mysql"
	"tiny_talk/net/router"
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

	// DBClient 初始化
	dsn := utils.ParseToDsn(&appConfig.Mysql)
	logger.Infof("dsn = %v", dsn)
	err = mysql.NewDBClient(dsn)
	if err != nil {
		logger.Fatalf("failed to init mysql, err = %v", err)
	}
	r := router.Router()
	r.Run(":8081")
}
