package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
)

// LogWithData data 可以是map, 也可以是一个结构体的指针
func LogWithData(module string, data any) *logrus.Entry {
	return logrus.WithField(constant.LogModule, module).WithField(constant.LogData, data)
}

func Log(module string) *logrus.Entry {
	return logrus.WithField(constant.LogModule, module)
}
func LogWithRequestId(ctx *gin.Context, module string, err error) *logrus.Entry {
	requestId, _ := ctx.Get(constant.RequestId)
	return logrus.WithField(constant.RequestId, requestId).WithField(constant.LogModule, module).WithError(err)
}

func LogWithRID(module string, data any, ctx *gin.Context) *logrus.Entry {
	requestId, _ := ctx.Get(constant.RequestId)
	return logrus.WithField(constant.RequestId, requestId).WithField(constant.LogModule, module).WithField(constant.LogData, data)
}

// LogBizErr 业务错误, 不需要module, 但是需要 errno.Errno
func LogBizErr(ctx *gin.Context, err errno.Errno, status int, info string) {
	requestId, _ := ctx.Get(constant.RequestId)
	logrus.WithField(constant.RequestId, requestId).WithError(err).Debug(info)
	ctx.AbortWithStatusJSON(status, err)
}

// LogParamError 特别针对参数错误
func LogParamError(ctx *gin.Context, err error) {
	requestId, _ := ctx.Get(constant.RequestId)
	errN := errno.Param.WithMessage(err.Error())
	ctx.AbortWithStatusJSON(http.StatusOK, errN)
	logrus.WithField(constant.RequestId, requestId).WithError(errN).Debug("request param err")
}

// LogParamErr 特别针对参数错误
func LogParamErr(ctx *gin.Context, err string) {
	requestId, _ := ctx.Get(constant.RequestId)
	errN := errno.Param.WithMessage(err)
	ctx.AbortWithStatusJSON(http.StatusOK, errN)
	logrus.WithField(constant.RequestId, requestId).WithError(errN).Debug("request param err")
}

// LogDB db err or dao return errno.Errno
func LogDB(ctx *gin.Context, err error) {
	requestId, _ := ctx.Get(constant.RequestId)
	if errors.Is(err, errno.Errno{}) {
		logrus.WithField(constant.RequestId, requestId).WithError(err).Debug()
	} else {
		logrus.WithField(constant.RequestId, requestId).WithError(errno.Service).Warn(err.Error())
	}
	ctx.AbortWithStatusJSON(http.StatusOK, err)
}
