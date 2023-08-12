package api

import (
	"github.com/gin-gonic/gin"
	"tiktok/service"
)

type UserLoginResp struct {
	StatusRespond
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
type UserResponse struct {
	StatusRespond
	User service.UserRespond `json:"user"`
}

func Register(c *gin.Context) {
}
func Login(c *gin.Context) {

}
func UserInfo(c *gin.Context) {}

func VideoList(c *gin.Context) {}
