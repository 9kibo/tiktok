package ginmw

import (
	"github.com/gin-gonic/gin"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
	"time"
)

type LogRequest struct {
	ClientIP       string
	Method         string
	Path           string
	ContentType    string
	BodySize       int
	StatusCode     int
	CompleteTimeMS time.Duration
}

// WithLogger gin的打印请求中间件
func WithLogger(skipPaths []string) gin.HandlerFunc {
	notlogged := skipPaths
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	return func(ctx *gin.Context) {
		requestId := utils.UUID4()
		ctx.Set(constant.RequestId, requestId)
		// Start timer
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		// Process request
		ctx.Next()
		if _, ok := skip[path]; ok {
			return
		}

		//no skip
		end := time.Now()
		param := LogRequest{}
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		param.ClientIP = ctx.ClientIP()
		param.Method = ctx.Request.Method
		if len(ctx.Request.Header["Content-Type"]) != 0 {
			param.ContentType = ctx.Request.Header["Content-Type"][0]
		}
		param.StatusCode = ctx.Writer.Status()
		param.CompleteTimeMS = end.Sub(start)
		param.BodySize = ctx.Writer.Size()
		utils.LogWithRequestIdData("requestLog", &param, ctx).Debug(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
	}
}
