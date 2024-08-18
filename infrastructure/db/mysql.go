package db

import (
	"tiny_talk/infrastructure/models"
	"tiny_talk/utils/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlClient *gorm.DB

func NewDBClient(dsn string) error {
	dbhandle, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger: gormlogger.New(
		// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		// 	gormlogger.Config{
		// 		SlowThreshold:             time.Second,     // 慢 SQL 阈值
		// 		LogLevel:                  gormlogger.Info, // 日志级别
		// 		IgnoreRecordNotFoundError: true,            // 忽略记录未找到的错误
		// 		Colorful:                  true,            // 彩色输出
		// 	},
		// ),
	})
	if err != nil {
		logger.Errorf("failed to connect database, err = %v", err)
		return err
	}
	if err := dbhandle.AutoMigrate(
		&models.UserBasic{},
		&models.FriendBasic{},
		&models.MessageBasic{},
		&models.ReceiverBox{}); err != nil {
		logger.Errorf("failed to migrate database, err = %v", err)
		return err
	}
	MysqlClient = dbhandle
	return nil
}
