package logmw

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"os"
	"tiktok/biz/config"
)

// InitLog
// 如果是文件, 需要使用日志分割功能, 可以使用软件日志分割功能, 比如Linux logrotate软件
func InitLog() (logger.LogLevel, logger.Writer) {
	loggerConfig := config.C.Logger
	//logrus
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(loggerConfig.GetLogrusLogLevel())
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     logrus.IsLevelEnabled(logrus.DebugLevel),
	})

	//gin
	ginMode := loggerConfig.GetGinMode()
	if ginMode != "" {
		gin.SetMode(ginMode)
	}
	redirectLogFromGin := NewRedirectLog(Gin)
	gin.DefaultWriter = redirectLogFromGin
	gin.DefaultErrorWriter = redirectLogFromGin

	//gorm
	return loggerConfig.GetGormLogLevel(), NewRedirectLog(Gorm)
}
