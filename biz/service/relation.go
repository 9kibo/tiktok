package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/redis"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

type FollowService interface {
	//FollowAction 关注或者取消关注
	FollowAction(req *model.FollowActionReq)
	// GetFollowingList   获取当前用户的关注列表
	GetFollowingList(userId int64) []*model.User
	// GetFollowerList 获取当前用户的粉丝列表
	GetFollowerList(userId int64) []*model.User
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
	if !isTokenUser(s.ctx, follow.FollowerId) {
		return
	}

	exists, err := dao.ExistsFollow(follow)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return
	}
	if req.ActionType == constant.FollowADD {
		//已关注不能关注
		if exists {
			utils.LogBizErr(s.ctx, errno.FollowAlreadyExist, http.StatusOK, "maybe a attack or front err")
			return
		}
		err = dao.AddFollow(follow)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return
		}
	} else {
		//未关注不能取消关注
		if !exists {
			utils.LogBizErr(s.ctx, errno.FollowNotExist, http.StatusOK, "has not follow,  maybe a attack or front err")
			return
		}
		err = dao.DeleteFollow(follow)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return
		}
	}
	redis.NewFollowService(s.ctx).DeleteFollowingIds(follow.FollowerId)
	redis.NewFollowService(s.ctx).DeleteFollowerIds(follow.FolloweeId)
}

func (s FollowServiceImpl) GetFollowingList(userId int64) []*model.User {

	//从缓存拿
	ids, ok := redis.NewFollowService(s.ctx).GetFollowingIds(userId)
	if !ok {
		//缓存没有从数据库拿, 再进行缓存
		followList, err := dao.GetFollowList(userId)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return nil
		}
		ids = s.getFollowingIds(followList)
		redis.NewFollowService(s.ctx).AddFollowingIds(userId, ids)
	}

	//根据ids查users返回
	users, err := dao.MustGetUsersByIds(ids)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return nil
	}
	return users
}
func (s FollowServiceImpl) getFollowingIds(follows []*model.Follow) []int64 {
	ids := make([]int64, 0, len(follows))
	for _, follow := range follows {
		ids = append(ids, follow.FolloweeId)
	}
	return ids
}
func (s FollowServiceImpl) GetFollowerList(userId int64) []*model.User {

	//从缓存拿
	ids, ok := redis.NewFollowService(s.ctx).GetFollowerIds(userId)
	if !ok {
		//缓存没有从数据库拿, 再进行缓存
		followList, err := dao.GetFollowerList(userId)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return nil
		}
		ids = s.getFollowerIds(followList)
		redis.NewFollowService(s.ctx).AddFollowerIds(userId, ids)
	}

	//根据ids查users返回
	users, err := dao.MustGetUsersByIds(ids)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return nil
	}
	return users
}
func (s FollowServiceImpl) getFollowerIds(follows []*model.Follow) []int64 {
	ids := make([]int64, 0, len(follows))
	for _, follow := range follows {
		ids = append(ids, follow.FollowerId)
	}
	return ids
}
