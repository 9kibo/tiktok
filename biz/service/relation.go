package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/logmw"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
)

type FollowService interface {
	//FollowAction 关注或者取消关注
	FollowAction(req *model.FollowActionReq)
	// GetFollowerList 获取当前用户的粉丝列表
	GetFollowerList(req *model.FollowingListReq) []*model.User
	// GetFolloweeList   获取当前用户的关注列表
	GetFolloweeList(userId int64) []*model.User
}
type FollowServiceImpl struct {
	//比如c.Get(constant.UserId)
	ctx *gin.Context
}

func NewFollowService(c *gin.Context) FollowService {
	return &FollowServiceImpl{
		ctx: c,
	}
}

func (s FollowServiceImpl) FollowAction(req *model.FollowActionReq) {
	follow := req.GetFollow()
	exists, err := dao.ExistsFollow(follow)
	if err != nil {
		s.ctx.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr)
		logmw.LogWithRequestId("follow", s.ctx).WithError(err).Debug("数据库异常")
		return
	}
	if req.ActionType == constant.FollowADD {
		//已关注不能关注
		if exists {
			s.ctx.AbortWithStatusJSON(http.StatusConflict, errno.FollowRelationAlreadyExistErr)
			logmw.LogWithRequestId("follow", s.ctx).Debug("已关注却重新关注, 不是攻击就是前端错误")
			return
		}
		err = dao.AddFollow(follow)
		if err != nil {
			s.ctx.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr)
			logmw.LogWithRequestId("follow", s.ctx).WithError(err).Debug("数据库异常")
			return
		}
	} else {
		//未关注不能取消关注
		if !exists {
			s.ctx.AbortWithStatusJSON(http.StatusConflict, errno.FollowRelationAlreadyExistErr)
			logmw.LogWithRequestId("follow", s.ctx).Debug("未关注却取消关注, 不是攻击就是前端错误")
			return
		}
		err = dao.DeleteFollow(follow)
		if err != nil {
			s.ctx.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr)
			logmw.LogWithRequestId("follow", s.ctx).WithError(err).Debug("数据库异常")
			return
		}
	}
}

func (s FollowServiceImpl) GetFollowerList(req *model.FollowingListReq) []*model.User {
	//followerList, err := dao.GetFollowerList(&model.Follow{
	//	FolloweeId: req.UserId,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//followerIdList := make([]int64, 0, len(followerList))
	//for _, follow := range followerList {
	//	followerIdList = append(followerIdList, follow.FollowerId)
	//}
	//TODO implement me
	panic("implement me")
}

func (s FollowServiceImpl) GetFolloweeList(userId int64) []*model.User {
	//TODO implement me
	panic("implement me")
}
