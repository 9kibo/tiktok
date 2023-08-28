package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/model"
	"tiktok/biz/service"
	"tiktok/pkg/constant"
)

type FavoriteActionResp struct {
	*model.BaseResp
}
type FavoriteListResp struct {
	*model.BaseResp
	VideoList []*model.Video `json:"video_list"`
}

func Favorite(c *gin.Context) {
	req := &model.FavoriteReq{}
	req.VideoId = c.GetInt64(constant.UserId)
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, model.BuildBindResp(err))
	}
	F := service.NewFavorite(c)
	F.FavouriteAction(req.UserId, req.VideoId, req.ActionType)
	if c.IsAborted() {
		return
	}
	respOfUpdate(c)

}

func FavoriteList(c *gin.Context) {
	req := &model.FavoriteListReq{}
	req.CurUserId = c.GetInt64(constant.UserId)
	if err := c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, model.BuildBindResp(err))
	}
	F := service.NewFavorite(c)
	List := F.GetFavouriteList(req.UserId, req.CurUserId)
	if c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, FavoriteListResp{
		model.BuildBaseResp(nil),
		List,
	})
}
