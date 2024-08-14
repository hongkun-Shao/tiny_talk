package main

import (
	"context"
	"log"
	"tiny_talk/infrastructure/db"
	"tiny_talk/net/router"
	"tiny_talk/utils"
	"tiny_talk/utils/config"
	"tiny_talk/utils/logger"
)

func main() {
	ctx := context.Background()
	log.Printf("tiny_talk start ...")
	path := utils.GetConfigPath()
	appConfig, err := config.LoadAppStaticConfig(path)
	if err != nil {
		log.Fatalf("failed to load config file, err = %v", err)
	}

	// 日志初始化
	_, err = logger.InitLogger(appConfig.Loggers, "text", "main")
	if err != nil {
		log.Panicf("failed to init text format logger, err = %v", err)
	}
	logger.Info("init logger success")

	// DBClient 初始化
	dsn := utils.ParseToDsn(&appConfig.Mysql)
	logger.Infof("dsn = %v", dsn)
	err = db.NewDBClient(dsn)
	if err != nil {
		logger.Panicf("failed to init mysql, err = %v", err)
	}
	logger.Info("init mysql success")

	// RDBClient 初始化
	err = db.NewRedisClient(&ctx, &appConfig.Redis)
	if err != nil {
		logger.Panicf("failed to init redis, err = %v", err)
	}
	logger.Info("init redis success")

	r := router.Router()
	r.Run(":8081")
}
