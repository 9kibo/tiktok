package mswagger

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	//执行swagger生成的go文件的初始化
	_ "tiktok/docs"
)

func InitSwagger(e *gin.Engine) {
	// 注册Swagger api相关路由
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
