package main

import (
	"tiktok/model"
	"tiktok/router"
)

func main() {
	model.InitDb()
	router.InitRouter()

}
