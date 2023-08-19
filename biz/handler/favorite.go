package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/model"
	"tiktok/biz/service"
)

type FavoriteActionResp struct {
	model.BaseResp
}
type FavoriteListResp struct {
	model.BaseResp
	//VideoList []service.VideoRespond `json:"video_list"`
}

func Favorite(c *gin.Context) {
	req := &model.FavoriteReq{}
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, model.BuildBindResp(err))
	}
	F := service.NewFavorite(c)
	if err := F.FavouriteAction(req.UserId, req.VideoId, req.ActionType); err != nil {
		return
	}
}

func FavoriteList(c *gin.Context) {}
