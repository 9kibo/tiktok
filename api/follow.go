package api

import (
	"github.com/gin-gonic/gin"
	"tiktok/service"
)

type RelationActionResp struct {
	StatusRespond
}

type FollowingListResp struct {
	StatusRespond
	UserList []service.UserRespond `json:"user_list"`
}
type FollowersListResp struct {
	StatusRespond
	UserList []service.UserRespond `json:"user_list"`
}

func RelationAction(c *gin.Context)   {}
func GetFollowingList(c *gin.Context) {}
func GetFollowersList(c *gin.Context) {}
func GetFriendList(c *gin.Context)    {}
