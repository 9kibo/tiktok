package dao

import (
	"tiktok/biz/config"
	"tiktok/pkg/log"
)

func init() {
	config.UseProfile("dev")
	config.Init("../../config.ini")
	gormLogLevel, gormLogWriter := log.InitLog()
	Init(gormLogLevel, gormLogWriter)
}
