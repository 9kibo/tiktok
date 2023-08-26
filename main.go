package main

import (
	"github.com/gin-gonic/gin"
	"tiktok/biz/config"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/ginmw"
	"tiktok/biz/middleware/kafka"
	"tiktok/biz/middleware/redis"
	"tiktok/pkg/log"
	"tiktok/pkg/swagger"
	"tiktok/pkg/validate"
)

func Init() {
	config.Init("config.ini")
	gormLogLevel, gormLogWriter := log.InitLog()
	dao.Init(gormLogLevel, gormLogWriter)
	kafka.Init()
	redis.Init()
	validate.InitValidateWrapper("")
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
// @BasePath /
// @query.collection.format multi
func main() {
	Init()
	e := gin.New()
	e.Use(ginmw.WithRecovery(), ginmw.WithLogger(nil))
	swagger.InitSwagger(e)
	initRouter(e)
	if err := e.Run(config.C.Server.Addr); err != nil {
		panic(err)
	}
}
