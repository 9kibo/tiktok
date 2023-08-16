package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"strings"
	"tiktok/pkg/constant"
)

var (
	C *Configuration
)

func Init(path string) {
	var err error
	profile := os.Getenv(constant.ProfileEnv)
	if profile != "" {
		log.Println("profile = ", profile)
		filename := filepath.Base(path)
		split := strings.Split(filename, ".")
		path = filepath.Join(filepath.Dir(path), split[0]+"-"+profile+"."+split[1])
	}

	C = &Configuration{}
	if err = ini.MapTo(C, path); err != nil {
		panic(fmt.Sprintf("配置文件读取错误, 请检查文件路径, err=%s", err))
	}
	C.Server.Profile = profile
}

// UseDevProfile 使用dev环境, 设置环境变量constant.ProfileEnv
func UseDevProfile() {
	os.Setenv(constant.ProfileEnv, "dev")
}

type Configuration struct {
	Profile string
	Server  *Server
	Logger  *Logger
	Jwt     *Jwt
	Mysql   *Mysql
	Cos     *Cos
	Redis   *Redis
	Kafka   *Kafka
}
type Server struct {
	Profile string
	Addr    string
}
type Logger struct {
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

type Jwt struct {
	SecretKey string
	Alg       string
	TokenKey  string
	Issuer    string
	Audience  string
	ExpireDay int64
}
type Mysql struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}
type Cos struct {
	Url       string
	ImageUrl  string
	SecretId  string
	SecretKey string
}
type Redis struct {
	Addr     string
	Password string
}
type Kafka struct {
	Addr string
}
