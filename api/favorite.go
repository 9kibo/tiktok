package api

import (
	"github.com/gin-gonic/gin"
	"tiktok/service"
)

type FavoriteActionResp struct {
	StatusRespond
}
type FavoriteListResp struct {
	StatusRespond
	VideoList []service.VideoRespond `json:"video_list"`
}

func Favorite(c *gin.Context)     {}
func FavoriteList(c *gin.Context) {}
