package main

import (
	"tiktok/model"
	"tiktok/router"
	"tiktok/utils/kafka"
	"tiktok/utils/redis"
)

func main() {
	model.InitDb()
	redis.InitRedis()
	kafka.InitKafka()
	router.InitRouter()

}
