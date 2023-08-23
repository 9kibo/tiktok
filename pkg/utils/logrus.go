package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"tiktok/pkg/constant"
)

// LogWithData data 可以是map, 也可以是一个结构体的指针
func LogWithData(module string, data any) *logrus.Entry {
	return logrus.WithField(constant.LogModule, module).WithField(constant.LogData, data)
}

func Log(module string) *logrus.Entry {
	return logrus.WithField(constant.LogModule, module)
}
func LogWithRequestId(module string, ctx *gin.Context) *logrus.Entry {
	requestId, _ := ctx.Get(constant.RequestId)
	return logrus.WithField(constant.LogModule, module).WithField(constant.RequestId, requestId)
}

func LogWithRequestIdData(module string, data any, ctx *gin.Context) *logrus.Entry {
	requestId, _ := ctx.Get(constant.RequestId)
	return logrus.WithField(constant.LogModule, module).WithField(constant.LogData, data).WithField(constant.RequestId, requestId)
}
