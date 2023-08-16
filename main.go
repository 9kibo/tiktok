package main

import (
	"github.com/gin-gonic/gin"
	"tiktok/biz/config"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/logmw"
	"tiktok/biz/middleware/mswagger"
)

func Init() {
	config.Init("config.ini")
	gormLogLevel, gormLogWriter := logmw.InitLog()
	dao.Init(gormLogLevel, gormLogWriter)
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
	e := gin.Default()
	mswagger.InitSwagger(e)
	initRouter(e)
	if err := e.Run(config.C.Server.Addr); err != nil {
		panic(err)
	}
}
