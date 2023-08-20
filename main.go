package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"tiktok/biz/config"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/kafka"
	"tiktok/biz/middleware/logmw"
	"tiktok/biz/middleware/mswagger"
)

func Init() {
	//logrus未初始化时使用控制台格式化json日志
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})
	config.Init("config.ini")
	gormLogLevel, gormLogWriter := logmw.InitLog()
	dao.Init(gormLogLevel, gormLogWriter)
	kafka.Init()
}

// @title mock tiktok
// @version 1.0 版本
// @description 字节青训营-模仿抖音项目
// @termsOfService http://swagger.io/terms/
// @contact.name 联系人
// @contact.url http://www.swagger.io/support
// @contact.email 584807419@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host videotools.cn
// @BasePath /
// @query.collection.format multi
func main() {
	Init()
	e := gin.New()
	e.Use(logmw.WithRecovery(), logmw.WithLogger(nil))
	mswagger.InitSwagger(e)
	initRouter(e)
	if err := e.Run(config.C.Server.Addr); err != nil {
		panic(err)
	}
}
