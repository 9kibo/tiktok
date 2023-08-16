package handler

import (
	"github.com/gin-gonic/gin"
	"tiktok/biz/model"
)

type FeedResp struct {
	model.BaseResp
	//VideoList []service.VideoRespond `json:"video_list"`
	NextTime int64 `json:"next_time"`
}

// Feed 视频Feed流：支持所有用户刷抖音，视频按投稿时间倒序推出
func Feed(c *gin.Context) {}

// UpVideo 视频投稿：支持登录用户自己拍视频投稿
func UpVideo(c *gin.Context)   {}
func VideoList(c *gin.Context) {}
