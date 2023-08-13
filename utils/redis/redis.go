package redis

import (
	"context"
)

var Ctx = context.Background()

/*
Redis连接在此声明
用户模块使用 db 0，1，2，3
视频模块使用 db 4，5，6，7
互动模块使用 db 8，9，10，11
社交模块使用 db 12，13，14，15
*/
//var ExampleRedis *redis.Client

func InitRedis() {
	/*
		示例
		ExampleRedis = redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPwd,
			DB:       15,
		})
	*/
}
