// Package redis
// get:
// add: need check success size, but set no for maybe exists element
// delete:  ignore if exists,  no Err as success delete
package redis

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"tiktok/pkg/utils"
)

const (
	follower  = "follower:"
	following = "following:"
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
func (s *FollowService) AddFollowingIds(userId int64, followingIds []int64) error {
	return sAdd(s.client, s.ctx, getIntKey(following, userId), expireS, utils.I2I(followingIds))
}
func (s *FollowService) AddFollowerIds(userId int64, followerIds []int64) error {
	return sAdd(s.client, s.ctx, getIntKey(follower, userId), expireS, utils.I2I(followerIds))
}

func (s *FollowService) DeleteFollowingIds(userId int64, ids []int64) error {
	var remCmd *redis.IntCmd
	if ids == nil {
		remCmd = s.client.SRem(s.ctx, getIntKey(following, userId))
	} else {
		remCmd = s.client.SRem(s.ctx, getIntKey(following, userId), utils.I2I(ids))
	}
	if remCmd.Err() != nil {
		return newErr(remCmd, "SRem")
	}
	return nil
}
func (s *FollowService) DeleteFollowerIds(userId int64, ids []int64) error {
	var remCmd *redis.IntCmd
	if ids == nil {
		remCmd = s.client.SRem(s.ctx, getIntKey(follower, userId))
	} else {
		remCmd = s.client.SRem(s.ctx, getIntKey(follower, userId), utils.I2I(ids))
	}
	if remCmd.Err() != nil {
		return newErr(remCmd, "SRem")
	}
	return nil
}

// GetFollowingIds 获取关注的人
func (s *FollowService) GetFollowingIds(userId int64) ([]int64, error) {
	return sGetInts(s.client, s.ctx, getIntKey(following, userId))
}

// GetFollowerIds 获取粉丝
func (s *FollowService) GetFollowerIds(userId int64) ([]int64, error) {
	return sGetInts(s.client, s.ctx, getIntKey(follower, userId))
}

func (s *FollowService) GetFollowingCount(userId int64) (int64, error) {
	return scard(s.client, s.ctx, getIntKey(following, userId))
}

func (s *FollowService) GetFollowerCount(userId int64) (int64, error) {
	return scard(s.client, s.ctx, getIntKey(follower, userId))
}

func (s *FollowService) ExistsFollow(followerId int64, followeeId int64) (bool, error) {
	return sisMember(s.client, s.ctx, following+strconv.FormatInt(followerId, 10), followeeId)
}
