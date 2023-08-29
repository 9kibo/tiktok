package redis

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
	"time"
)

func sGetInts(client redis.Cmdable, ctx *gin.Context, key string) ([]int64, error) {
	cmd := client.SMembers(ctx, key)
	if cmd.Err() != nil {
		return nil, newErr(cmd, "SMembers")
	}
	vals := cmd.Val()
	ints := make([]int64, 0, len(vals))
	for _, v := range vals {
		parseInt, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, newErr(cmd, fmt.Sprintf("%s ' element must a int64, elements=[%v]", key, vals))
		}
		ints = append(ints, parseInt)
	}
	return ints, nil
}
func scard(client redis.Cmdable, ctx *gin.Context, key string) (int64, error) {
	cmd := client.SCard(ctx, key)
	if cmd.Err() != nil {
		return 0, newErr(cmd, "scard")
	}
	return cmd.Val(), nil
}
func sisMember(client redis.Cmdable, ctx *gin.Context, key string, member any) (bool, error) {
	cmd := client.SIsMember(ctx, key, member)
	if cmd.Err() != nil {
		return false, newErr(cmd, "SIsMember")
	}
	return cmd.Val(), nil
}
func sAdd(client redis.Cmdable, ctx *gin.Context, key string, expireD time.Duration, values []any) error {
	//Err is Pipe func return's Err or redis request Err
	cmds, err := client.TxPipelined(ctx, func(p redis.Pipeliner) error {
		sAddCmd := p.SAdd(ctx, key, values...)
		if sAddCmd.Err() != nil {
			return sAddCmd.Err()
		} else if sAddCmd.Val() != int64(len(values)) {
			return fmt.Errorf("redis sadd wrong size, sadd %d, but result=%d", sAddCmd.Val(), len(values))
		}
		expireCmd := p.Expire(ctx, key, expireD)
		if expireCmd.Err() != nil {
			return expireCmd.Err()
		}
		return nil
	})
	if err != nil {
		utils.LogWithRequestId(ctx, constant.LMRedis, err).Debug("Cmd=" + cmdsString(cmds))
		return newErrCmds(cmds, "SAdd", err)
	}
	return nil
}
