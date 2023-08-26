package redis

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
)

func sGetInt(client redis.Cmdable, ctx *gin.Context, key string) ([]int64, bool) {
	cmd := client.SMembers(context.Background(), key)
	if handlerErr(cmd, ctx) {
		return nil, false
	}
	vals := cmd.Val()
	ints := make([]int64, 0, len(vals))
	for _, v := range vals {
		parseInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			utils.Log(constant.LMRedis).WithError(err).Debugf("%s ' element must a int64, elements=[%v]", key, vals)
			return nil, false
		}
		ints = append(ints, parseInt)
	}
	return ints, true
}
func scard(client redis.Cmdable, ctx *gin.Context, key string) (int64, bool) {
	cmd := client.SCard(context.Background(), key)
	if handlerErr(cmd, ctx) {
		return 0, false
	}
	return cmd.Val(), true
}
func sisMember(client redis.Cmdable, ctx *gin.Context, key string, member any) (bool, bool) {
	cmd := client.SIsMember(context.Background(), key, member)
	if handlerErr(cmd, ctx) {
		return false, false
	}
	return cmd.Val(), true
}
