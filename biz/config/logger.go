package config

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"os"
)

type Logger struct {
	Filepath    string
	LogrusLevel string
	GormLevel   string
	GinLevel    string
}

func (l Logger) GetGinMode() string {
	if os.Getenv(gin.EnvGinMode) != "" {
		return ""
	}
	switch l.GinLevel {
	case "DEBUG":
		return gin.DebugMode
	case "RELEASE":
		return gin.ReleaseMode
	case "TEST":
		return gin.TestMode
	}
	panic("has not GinLevel=" + l.GinLevel)
}
func (l Logger) GetGormLogLevel() logger.LogLevel {
	switch l.GormLevel {
	case "ERROR":
		return logger.Error
	case "WARN":
		return logger.Warn
	case "INFO":
		return logger.Info
	case "SILENT":
		return logger.Silent
	}
	panic("has not GormLevel=" + l.GormLevel)
}
func (l Logger) GetLogrusLogLevel() logrus.Level {
	switch l.LogrusLevel {
	case "PANIC":
		return logrus.PanicLevel
	case "FATAL":
		return logrus.FatalLevel
	case "ERROR":
		return logrus.ErrorLevel
	case "WARN":
		return logrus.WarnLevel
	case "INFO":
		return logrus.InfoLevel
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	}
	panic("has not LogrusLevel=" + l.LogrusLevel)
}
