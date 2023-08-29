package redis

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"tiktok/pkg/utils"
)

const (
	friendList = "friendList:"
)

type MessageService struct {
	client *redis.Client
	ctx    *gin.Context
}

func NewMessageService(ctx *gin.Context) *MessageService {
	return &MessageService{
		client: messageClient,
		ctx:    ctx,
	}
}

func (s *MessageService) AddFriendIds(userId int64, friendIds []int64) error {
	return sAdd(s.client, s.ctx, getIntKey(friendList, userId), expireS, utils.I2I(friendIds))
}

func (s *MessageService) DeleteFriendIds(userId int64) error {
	remCmd := s.client.SRem(s.ctx, getIntKey(friendList, userId))
	if remCmd.Err() != nil {
		return newErr(remCmd, "SRem")
	}
	return nil
}

// GetFriendIds
// @return FriendIds and
func (s *MessageService) GetFriendIds(userId int64) ([]int64, error) {
	return sGetInts(s.client, s.ctx, getIntKey(friendList, userId))
}
