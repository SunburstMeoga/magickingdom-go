package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"magickingdom-go/internal/config"
)

var Log *logrus.Logger

// Init 初始化日志
func Init(cfg *config.LogConfig) error {
	Log = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	Log.SetLevel(level)

	// 设置日志格式
	if cfg.Format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// 设置输出
	if cfg.Output == "file" {
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		Log.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		Log.SetOutput(os.Stdout)
	}

	return nil
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	if Log == nil {
		Log = logrus.New()
	}
	return Log
}

