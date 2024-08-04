package mysql

import (
	"log"
	"os"
	"time"
	"tiny_talk/infrastructure/mysql/models"
	"tiny_talk/utils/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type DBClient struct {
	dbhandle *gorm.DB
}

var MysqlClient *DBClient

func NewDBClient(dsn string) error {
	dbhandle, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			gormlogger.Config{
				SlowThreshold:             time.Second,     // 慢 SQL 阈值
				LogLevel:                  gormlogger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,            // 忽略记录未找到的错误
				Colorful:                  true,            // 彩色输出
			},
		),
	})
	if err != nil {
		logger.Errorf("failed to connect database, err = %v", err)
		return err
	}
	if err := dbhandle.AutoMigrate(&models.UserBasic{}); err != nil {
		logger.Errorf("failed to migrate databaseerr = %v", err)
		return err
	}
	MysqlClient = &DBClient{
		dbhandle: dbhandle,
	}
	return nil
}

type ModelInterface interface {
	TableName() string
}

type ModelCRUD interface {
	ModelInterface
	Create(db *gorm.DB) error
	Update(db *gorm.DB) error
	Find(db *gorm.DB) error
	Delete(db *gorm.DB) error
}

func CreateModel(model ModelCRUD) error {
	return model.Create(MysqlClient.dbhandle)
}

func UpdateModel(model ModelCRUD) error {
	return model.Update(MysqlClient.dbhandle)
}

func FindModel(model ModelCRUD) error {
	return model.Find(MysqlClient.dbhandle)
}

func DeleteModel(model ModelCRUD) error {
	return model.Delete(MysqlClient.dbhandle)
}
