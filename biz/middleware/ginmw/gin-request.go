package ginmw

import (
	"github.com/gin-gonic/gin"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
	"time"
)

type LogRequest struct {
	ClientIP     string
	Method       string
	Path         string
	RespBodySize int
	StatusCode   int
	LatencyMs    float64
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
		param.StatusCode = ctx.Writer.Status()
		latency := end.Sub(start)
		if latency < time.Millisecond {
			latency = latency.Truncate(time.Nanosecond)
			param.LatencyMs = float64(latency.Microseconds()) / 1000
		} else {
			param.LatencyMs = float64(latency.Milliseconds())
		}
		param.RespBodySize = ctx.Writer.Size()
		utils.LogWithRID("requestLog", &param, ctx).Debug(ctx.Errors.ByType(gin.ErrorTypePrivate).String())
	}
}
