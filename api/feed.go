package api

import (
	"github.com/gin-gonic/gin"
	"tiktok/service"
)

type FeedResp struct {
	StatusRespond
	VideoList []service.VideoRespond `json:"video_list"`
	NextTime  int64                  `json:"next_time"`
}

func Feed(c *gin.Context) {}

func UpVideo(c *gin.Context) {}
