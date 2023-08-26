package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"os"
)

var (
	C *Configuration
)

func init() {
	//logrus未初始化时使用控制台格式化json日志
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})
}
func Init(path string) {
	var err error

	profile := GetActiveProfile()
	profilePath := ""
	if profile != "" {
		profilePath = GetProfilePath(path, profile)
	}
	c := Configuration{}
	if profile != "" {
		//环境变量的profile
		if err = ini.MapTo(&c, path, profilePath); err != nil {
			panic(fmt.Sprintf("配置文件读取错误, 请检查文件路径, err=%s", err))
		}
		c.Server.Profile = profile
	} else {
		if err = ini.MapTo(&c, path); err != nil {
			panic(fmt.Sprintf("配置文件读取错误, 请检查文件路径, err=%s", err))
		}
		//配置文件中的profile
		if c.Server.Profile != "" {
			profilePath = GetProfilePath(path, profile)
			if err = ini.MapTo(&c, path, profilePath); err != nil {
				panic(fmt.Sprintf("配置文件读取错误, 请检查文件路径, err=%s", err))
			}
		}
	}
	C = &c
}

type Configuration struct {
	Server *Server
	Logger *Logger
	Jwt    *Jwt
	Mysql  *Mysql
	Cos    *Cos
	Redis  *Redis
	Kafka  *Kafka
}
type Server struct {
	Profile string
	Addr    string
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
