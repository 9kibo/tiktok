package dao

import (
	"os"
	"testing"
	"tiktok/biz/config"
	"tiktok/biz/middleware/logmw"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
)

func init() {
	os.Setenv(constant.ProfileEnv, "dev")
	config.Init("../../config.ini")
	gormLogLevel, gormLogWriter := logmw.InitLog()
	Init(gormLogLevel, gormLogWriter)
}
func TestA(t *testing.T) {
	Db.Find(&model.User{}).Where("id=123")
}
