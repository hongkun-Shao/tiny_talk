package logger

import (
	"log"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger = logrus.Logger
type Hook = logrus.Hook

var LoggerConf *logrus.Entry

type LoggerConfig struct {
	Level  string
	Output string
}

type LoggerMap struct {
	loggerConfigs map[string]LoggerConfig
	loggers       map[string]*Logger
}

var Loggers LoggerMap
var AppName string

func InitLogger(loggerConfig map[string]LoggerConfig, logType string, appName string) (*LoggerMap, error) {
	loggers := make(map[string]*Logger)

	for name, config := range loggerConfig {
		logger := logrus.New()
		log.Printf("logger %s init success", name)
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return nil, err
		}

		logger.SetLevel(level)

		logger.SetOutput(&lumberjack.Logger{
			Filename:   config.Output,
			MaxSize:    512,
			MaxBackups: 7,
			MaxAge:     7,
			Compress:   true,
		})
		if logType == "text" {
			logger.SetFormatter(generaterTextFormatter())
		} else {
			logger.SetFormatter(generaterJsonFormatter())
		}
		logger.SetReportCaller(true)
		loggers[name] = logger
	}
	AppName = appName
	Loggers = LoggerMap{
		loggerConfig,
		loggers,
	}
	return &Loggers, nil
}

func generaterTextFormatter() *nested.Formatter {
	return &nested.Formatter{
		CallerFirst:      true,
		HideKeys:         true,
		NoColors:         true,
		NoFieldsColors:   true,
		NoFieldsSpace:    true,
		NoUppercaseLevel: true,
		ShowFullLevel:    true,
		TimestampFormat:  "2006-01-02 15:04:05.000",
	}
}

func generaterJsonFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05.000",
		DisableHTMLEscape: true,
	}
}

func Debugf(format string, args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Infof(format, args...)
}

func Printf(format string, args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Printf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Panicf(format, args...)
}

func Debug(args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Debug(args...)
}

func Info(args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Info(args...)
}

func Print(args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Print(args...)
}

func Warn(args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Warn(args...)
}

func Error(args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Error(args...)
}

func Fatal(args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Fatal(args...)
}

func Panic(args ...interface{}) {
	Loggers.loggers[AppName].WithFields(logrus.Fields{}).Panic(args...)
}
