package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
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

	//GetFriendList 获取发生过消息的用户即好友, 顺便返回最新一条消息
	GetFriendList(userId int64) []*model.FriendUser
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
func (s FollowServiceImpl) GetFriendList(userId int64) []*model.FriendUser {
	if !isTokenUser(s.ctx, userId) {
		return nil
	}

	//get from redis
	redisService := redis.NewMessageService(s.ctx)
	friendIds, err := redisService.GetFriendIds(userId)
	if err != nil {
		redis.HandlerErr(s.ctx, err)
		return nil
	}

	//redis has not, get from db and add to redis
	friendIds, err = dao.GetFriendIds(userId)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return nil
	}
	//add to redis
	if err = redisService.AddFriendIds(userId, friendIds); err != nil {
		redis.HandlerErr(s.ctx, err)
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	var friends []*model.User
	var messages []*model.Message
	var fe, me error
	go func() {
		defer wg.Done()
		friends, err = dao.MustGetUsersByIds(friendIds)
		if err != nil {
			fe = err
			utils.LogDB(s.ctx, err)
		}
	}()
	go func() {
		defer wg.Done()
		messages, err = dao.GetMessageLatestByToUserIds(userId, friendIds)
		if err != nil {
			me = err
			utils.LogDB(s.ctx, err)
		}
	}()
	if fe != nil || me != nil {
		return nil
	}
	messageMap := make(map[int64]*model.Message, len(messages))
	for _, message := range messages {
		if message.FromUserId == userId {
			messageMap[message.ToUserId] = message
		} else {
			messageMap[message.FromUserId] = message
		}
	}
	friendList := make([]*model.FriendUser, 0, len(messages))
	for _, friend := range friends {
		var msgType int64
		latestMessage := messageMap[friend.Id]
		if latestMessage.FromUserId == userId {
			msgType = constant.MessageTypeSend
		} else {
			msgType = constant.MessageTypeAccept
		}
		friendList = append(friendList, &model.FriendUser{
			User:    friend,
			Message: messageMap[friend.Id].Content,
			MsgType: msgType,
		})
	}
	return friendList
}
func (s FollowServiceImpl) GetFollowingList(userId int64) []*model.User {

	//get from redis
	redisService := redis.NewFollowService(s.ctx)
	ids, err := redisService.GetFollowingIds(userId)
	if err != nil {
		redis.HandlerErr(s.ctx, err)
		return nil
	}
	if len(ids) == 0 {
		//redis has not, get from db and add to redis
		followList, err := dao.GetFollowList(userId)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return nil
		}
		ids = s.getFollowingIds(followList)
		//add to redis
		if err = redisService.AddFollowingIds(userId, ids); err != nil {
			redis.HandlerErr(s.ctx, err)
			return nil
		}
	}

	//get users
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

	//get from redis
	redisService := redis.NewFollowService(s.ctx)
	ids, err := redisService.GetFollowerIds(userId)
	if err != nil {
		redis.HandlerErr(s.ctx, err)
		return nil
	}
	//redis has not, get from db and add to redis
	if len(ids) == 0 {
		followList, err := dao.GetFollowerList(userId)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return nil
		}
		ids = s.getFollowerIds(followList)
		//add to redis
		err = redisService.AddFollowerIds(userId, ids)
		if err != nil {
			redis.HandlerErr(s.ctx, err)
			return nil
		}
	}

	//get users
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

	//delete
	redisService := redis.NewFollowService(s.ctx)
	err = redisService.DeleteFollowingIds(follow.FollowerId, nil)
	if err != nil {
		redis.HandlerErr(s.ctx, err)
		return
	}
	err = redisService.DeleteFollowerIds(follow.FolloweeId, nil)
	if err != nil {
		redis.HandlerErr(s.ctx, err)
		return
	}
}
