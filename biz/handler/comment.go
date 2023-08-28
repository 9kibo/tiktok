package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/model"
	"tiktok/biz/service"
	"tiktok/pkg/constant"
)

type CommentActionResp struct {
	*model.BaseResp
	Comment *model.Comment `json:"comment,omitempty"`
}
type CommentListResp struct {
	*model.BaseResp
	CommentList []*model.Comment `json:"comment_list"`
}

func CommAction(c *gin.Context) {
	var err error
	var CommResp *model.Comment
	req := &model.CommReq{}
	req.UserId = c.GetInt64(constant.UserId)
	if err = c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, model.BuildBindResp(err))
	}
	Comm := &service.CommentServiceImpl{C: c}
	if req.Action == 1 {
		CommResp, err = Comm.AddComment(req.UserId, req.VideoId, req.Text)
		if err != nil || c.IsAborted() {
			return
		}
	} else {
		err = Comm.DelComment(req.VideoId, req.DelCommId)
		if err != nil || c.IsAborted() {
			return
		}
	}
	c.JSON(http.StatusOK, CommentActionResp{
		BaseResp: model.BuildBaseResp(err),
		Comment:  CommResp,
	})
}

func CommList(c *gin.Context) {
	var err error
	req := &model.CommentsReq{}
	req.UserId = c.GetInt64(constant.UserId)
	if err = c.ShouldBindQuery(req); err != nil {
		c.JSON(http.StatusBadRequest, model.BuildBindResp(err))
	}
	Comm := &service.CommentServiceImpl{C: c}
	commentList, err := Comm.GetCommentList(req.VideoId)
	if err != nil || c.IsAborted() {
		return
	}
	c.JSON(http.StatusOK, CommentListResp{
		BaseResp:    model.BuildBaseResp(err),
		CommentList: commentList,
	})
}
