package handler

import (
	"github.com/gin-gonic/gin"
	"tiktok/biz/model"
)

type CommentActionResp struct {
	model.BaseResp
	//Comment service.CommentInfo `json:"comment,omitempty"`
}
type CommentListResp struct {
	model.BaseResp
	//CommentList []service.CommentInfo `json:"comment_list"`
}

func CommAction(c *gin.Context) {

}

func CommList(c *gin.Context) {}
