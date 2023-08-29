package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"tiktok/biz/config"
	"time"
)

var Ctx = context.Background()

/*
用户模块使用 db 0，1，2，3
视频模块使用 db 4，5，6，7
互动模块使用 db 8，9，10，11
社交模块使用 db 12，13，14，15
*/
var (
	expireS       = 60 * 60 * time.Second
	followClient  *redis.Client
	messageClient *redis.Client
)

func Init() {
	var err error
	followClient, err = GetRedis(12)
	if err != nil {
		panic(err)
	}

	messageClient, err = GetRedis(13)
	if err != nil {
		panic(err)
	}

}
func GetRedis(db int) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.Addr,
		Password: config.C.Redis.Password,
		DB:       db,
	})
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
