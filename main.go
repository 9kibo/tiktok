package main

import (
	"tiktok/model"
	"tiktok/router"
	"tiktok/utils/redis"
)

func main() {
	model.InitDb()
	redis.InitRedis()
	router.InitRouter()

}
