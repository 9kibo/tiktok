// Package redis
// get:
// add: need check success size, but set no for maybe exists element
// delete:  ignore if exists,  no err as success delete
package redis

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
	"time"
)

const (
	follower  = "follower:"
	following = "following:"
	expireS   = 60 * 60 * time.Second
)

type FollowService struct {
	client *redis.Client
	ctx    *gin.Context
}

func NewFollowService(ctx *gin.Context) *FollowService {
	return &FollowService{
		client: followClient,
		ctx:    ctx,
	}
}
func (s *FollowService) getKey(prefix string, userId int64) string {
	return prefix + strconv.FormatInt(userId, 10)
}
func (s *FollowService) AddFollowingIds(userId int64, followingIds []int64) bool {
	return s.addFollowIds(following, userId, followingIds)
}
func (s *FollowService) AddFollowerIds(userId int64, followerIds []int64) bool {
	return s.addFollowIds(follower, userId, followerIds)
}
func (s *FollowService) addFollowIds(keyPrefix string, userId int64, followingIds []int64) bool {
	//err is Pipe func return's err or request err
	cmds, err := s.client.TxPipelined(context.Background(), func(p redis.Pipeliner) error {
		key := s.getKey(keyPrefix, userId)
		size := len(followingIds)
		followingIds0 := make([]any, 0, size)
		for _, id := range followingIds {
			followingIds0 = append(followingIds0, id)
		}
		sAdd := p.SAdd(context.Background(), key, followingIds0...)
		if sAdd.Err() != nil {
			handlerErr(sAdd, s.ctx)
			return nil
		}
		expire := p.Expire(context.Background(), key, expireS)
		if expire.Err() != nil {
			handlerErr(sAdd, s.ctx)
			return nil
		}
		return nil
	})
	if err != nil {
		utils.LogWithRequestId(s.ctx, constant.LMRedis, err).Debug("cmd=%s", cmdsString(cmds))
		return false
	}
	return true
}

func (s *FollowService) DeleteFollowingIds(userId int64) bool {
	return !handlerErr(s.client.SRem(context.Background(), s.getKey(following, userId)), s.ctx)
}
func (s *FollowService) DeleteFollowerIds(userId int64) bool {
	return !handlerErr(s.client.SRem(context.Background(), s.getKey(follower, userId)), s.ctx)
}

// GetFollowingIds 获取关注的人
func (s *FollowService) GetFollowingIds(userId int64) ([]int64, bool) {
	return sGetInt(s.client, s.ctx, s.getKey(following, userId))
}

// GetFollowerIds 获取粉丝
func (s *FollowService) GetFollowerIds(userId int64) ([]int64, bool) {
	return sGetInt(s.client, s.ctx, s.getKey(follower, userId))
}

func (s *FollowService) GetFollowingCount(userId int64) (int64, bool) {
	return scard(s.client, s.ctx, s.getKey(following, userId))
}

func (s *FollowService) GetFollowerCount(userId int64) (int64, bool) {
	return scard(s.client, s.ctx, s.getKey(follower, userId))
}

func (s *FollowService) ExistsFollow(followerId int64, followeeId int64) (bool, bool) {
	return sisMember(s.client, s.ctx, following+strconv.FormatInt(followerId, 10), followeeId)
}
