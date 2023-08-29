package redis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"strings"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

type Err struct {
	Msg string
	Cmd string
	Err error
}

func (e Err) Error() string {
	return fmt.Sprintf("Msg=%s, Cmd=%s, redisErr=%s", e.Msg, e.Cmd, e.Err)
}

type logRedisCmdData struct {
	Cmd string
}

func newErr(cmd redis.Cmder, msg string) error {
	return &Err{
		Msg: msg,
		Cmd: cmd.String(),
		Err: cmd.Err(),
	}
}
func newErrCmds(cmd []redis.Cmder, msg string, err error) error {
	if redisErr, ok := err.(*Err); ok {
		return &Err{
			Msg: msg + ", " + redisErr.Msg,
			Cmd: redisErr.Cmd,
			Err: redisErr.Err,
		}
	}
	return &Err{
		Msg: msg,
		Cmd: cmdsString(cmd),
		Err: err,
	}
}

// HandlerErr 有错误就打印错误
func HandlerErr(ctx *gin.Context, err error) {
	redisErr := err.(*Err)
	utils.LogWithRID(constant.LMRedis, redisErr.Cmd, ctx).WithError(redisErr.Err).Debug(redisErr.Msg)
	ctx.AbortWithStatusJSON(http.StatusOK, errno.Service)
}

// 返回是否有错误并打印日志
func cmdsString(cmds []redis.Cmder) string {
	sb := strings.Builder{}
	size := len(cmds)
	for i, cmd := range cmds {
		sb.WriteString(strconv.FormatInt(int64(i+1), 10))
		sb.WriteString(": ")
		sb.WriteString(cmd.String())
		if i+1 != size {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

func getIntKey(prefix string, userId int64) string {
	return prefix + strconv.FormatInt(userId, 10)
}
