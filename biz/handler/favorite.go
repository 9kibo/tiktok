package handler

import (
	"github.com/gin-gonic/gin"
	"tiktok/biz/model"
)

type FavoriteActionResp struct {
	model.BaseResp
}
type FavoriteListResp struct {
	model.BaseResp
	//VideoList []service.VideoRespond `json:"video_list"`
}

func Favorite(c *gin.Context)     {}
func FavoriteList(c *gin.Context) {}
