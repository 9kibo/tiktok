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
	VideoList []model.Video `json:"video_list"`
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
	c.JSON(http.StatusOK, FavoriteActionResp{model.BaseResp{
		Code: 1,
		Msg:  "点赞成功",
	}})

}

func FavoriteList(c *gin.Context) {}
