package logmw

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	LogFieldName = "request"
	Module       = "module"
)

type LogRequest struct {
	ClientIP           string
	Method             string
	Path               string
	ContentType        string
	BodySize           int
	StatusCode         int
	CompleteTimeSecond time.Duration
}

// WithLoggerWithLogrus gin的打印请求中间件
func WithLoggerWithLogrus(skipPaths []string) gin.HandlerFunc {
	notlogged := skipPaths
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		c.Set("startRequestTime", start)
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()
		if c.IsAborted() {
			//自己处理了
			return
		}
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := GetRequestParams(c)
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			logrus.WithField(LogFieldName, &param).Debug(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
	}
}
func GetRequestParams(ctx *gin.Context) *LogRequest {
	end := time.Now()
	path := ctx.Request.URL.Path
	raw := ctx.Request.URL.RawQuery
	param := LogRequest{}
	param.ClientIP = ctx.ClientIP()
	param.Method = ctx.Request.Method
	if len(ctx.Request.Header["Content-Type"]) != 0 {
		param.ContentType = ctx.Request.Header["Content-Type"][0]
	}
	param.StatusCode = ctx.Writer.Status()

	if start, ok := ctx.Get("startRequestTime"); ok {
		param.CompleteTimeSecond = end.Sub(start.(time.Time))
	}

	param.BodySize = ctx.Writer.Size()
	if raw != "" {
		path = path + "?" + raw
	}
	param.Path = path
	return &param
}
func LogWithRequestErr(module string, ctx *gin.Context, err error) *logrus.Entry {
	params := GetRequestParams(ctx)
	return logrus.WithField(LogFieldName, params).WithField(Module, module).WithField("bizErr", err)
}
func LogWithRequest(module string, ctx *gin.Context) *logrus.Entry {
	params := GetRequestParams(ctx)
	return logrus.WithField(LogFieldName, params).WithField(Module, module)
}

const (
	Gin  = "GIN"
	Gorm = "GORM"
)

// RedirectLog
// 关于为什么不用Level
//
//	gin有gin.IsDebugging判断
//	gorm需动态判断
type RedirectLog struct {
	From string
}

func NewRedirectLog(from string) *RedirectLog {
	switch from {
	case Gorm, Gin:
	default:
		panic("NewRedirectLog: has the the from=" + from)
	}
	return &RedirectLog{
		From: from,
	}
}

// Write gin
func (w RedirectLog) Write(p []byte) (int, error) {
	if w.From == Gin {
		w.ginLog(string(p))
	}
	return len(p), nil
}

// Printf gorm
func (w RedirectLog) Printf(format string, args ...interface{}) {
	if w.From == Gorm {
		w.gormLog(format, args)
	}
}
func (w RedirectLog) ginLog(s string) {
	if gin.IsDebugging() {
		logrus.WithField(Module, w.From).Debug(strings.Replace(s, "[GIN-debug] ", "", 1))
	}
}

var (
	gormInfoStr = "%s\n[info] "
	gormWarnStr = "%s\n[warn] "
	gormErrStr  = "%s\n[error] "
)

// 必须禁止颜色
func (w RedirectLog) gormLog(format string, args ...interface{}) {
	var level logrus.Level
	if strings.HasPrefix(format, gormInfoStr) {
		level = logrus.InfoLevel
	} else if strings.HasPrefix(format, gormWarnStr) {
		level = logrus.WarnLevel
	} else if strings.HasPrefix(format, gormErrStr) {
		level = logrus.ErrorLevel
	} else {
		level = logrus.DebugLevel
	}

	log := logrus.WithField(Module, w.From).WithField("gormFile", args[0])
	if level == logrus.DebugLevel {
		if len(args) == 4 {
			//l.LogLevel == Info
			log.WithField("gormElapsed", args[1]).WithField("gormRows", args[2]).WithField("gormSql", args[3]).Debug("info")
		} else {
			log = log.WithField("gormElapsed", args[2]).WithField("gormRows", args[3]).WithField("gormSql", args[4])
			if _, ok := args[1].(string); ok {
				//l.LogLevel >= Warn
				log.WithField("gormSlowLog", args[1]).Debug("warn")
			} else {
				//l.LogLevel >= Error
				log.WithField("gormErr", args[1]).Debug("error")
			}
		}
	} else {
		log.WithField("gormData", args[1:]).Log(level)
	}
}
