package api

import (
	"github.com/gin-gonic/gin"
	"tiktok/service"
)

type CommentActionResp struct {
	StatusRespond
	Comment service.CommentInfo `json:"comment,omitempty"`
}
type CommentListResp struct {
	StatusRespond
	CommentList []service.CommentInfo `json:"comment_list"`
}

func CommAction(c *gin.Context) {
}

func CommList(c *gin.Context) {}
