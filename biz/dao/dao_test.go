package dao

import (
	"testing"
	"tiktok/biz/config"
	"tiktok/pkg/log"
)

func TestMain(m *testing.M) {
	config.UseProfile("dev")
	config.Init("../../config.ini")
	config.C.Logger.Filepath = "A:\\code\\backend\\go\\tiktok\\dao_test.log"
	gormLogLevel, gormLogWriter := log.InitLog()
	Init(gormLogLevel, gormLogWriter)
	m.Run()
}
