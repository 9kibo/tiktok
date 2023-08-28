package redis

import (
	"github.com/redis/go-redis/v9"
	"sync"
	"tiktok/pkg/utils"
)

// 复用全局连接
type CommRedis struct {
	Comment *redis.Client
	sync.Once
}

var CommMutex = utils.NewCacheGuard()

var CommR = FavRedis{}

func (Comm *CommRedis) GetCommRedis() (*redis.Client, error) {
	var err error
	Comm.Do(func() {
		client, err := GetRedis(9)
		if err != nil {
			utils.Log("redis").WithField("err:", err).Error("redis连接失败")
		}
		Comm.Comment = client
	})
	err = Comm.Comment.Ping(Ctx).Err()
	if err != nil {
		utils.Log("redis").WithField("err:", err).Error("redis连接失败")
		return nil, err
	}
	return Comm.Comment, nil
}
