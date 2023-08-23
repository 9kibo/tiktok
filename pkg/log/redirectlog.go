package log

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"strings"
	"tiktok/pkg/utils"
)

const (
	Gin         = "GIN"
	Gorm        = "GORM"
	GinRecovery = "GIN_RECOVERY"
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
	case Gorm, Gin, GinRecovery:
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
		utils.Log(w.From).Debug(strings.Replace(s, "[GIN-debug] ", "", 1))
	}
}

var (
	gormInfoStr = "%s\n[info] "
	gormWarnStr = "%s\n[warn] "
	gormErrStr  = "%s\n[error] "
)

// 必须禁止颜色
func (w RedirectLog) gormLog(format string, args []interface{}) {
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

	logData := make(map[string]any, 4)
	logData["file"] = args[0]
	if level == logrus.DebugLevel {
		if len(args) == 4 {
			//l.LogLevel == Info
			logData["elapsed"] = args[1]
			logData["rows"] = args[2]
			logData["sql"] = args[3]
			utils.LogWithData(w.From, logData).Debug("trace info")
		} else {
			logData["elapsed"] = args[2]
			logData["rows"] = args[3]
			logData["sql"] = args[4]
			if _, ok := args[1].(string); ok {
				//l.LogLevel >= Warn
				logData["slowLog"] = args[1]
				utils.LogWithData(w.From, logData).Debug("trace warn")
			} else {
				//l.LogLevel >= Error
				logData["err"] = args[1]
				utils.LogWithData(w.From, logData).Debug("trace error")
			}
		}
	} else {
		utils.LogWithData(w.From, args[1:]).Log(level)
	}
}
