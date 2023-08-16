package handler

import (
	"github.com/gin-gonic/gin"
	"tiktok/biz/model"
)

type UserRegisterReq struct {
	Username int64  `json:"username"`
	Password string `json:"password"`
}
type UserRegisterResp struct {
	model.BaseResp
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func Register(c *gin.Context) {

}

type UserLoginReq struct {
	Username int64  `json:"username"`
	Password string `json:"password"`
}
type UserLoginResp struct {
	model.BaseResp
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

func Login(c *gin.Context) {

}

// UserInfo 个人主页：支持查看用户基本信息和投稿列表，注册用户流程简化
func UserInfo(c *gin.Context) {

}
