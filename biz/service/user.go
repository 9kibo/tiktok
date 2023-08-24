package service

import (
	"github.com/gin-gonic/gin"
	"tiktok/biz/model"
)

type UserService interface {
	//GetUserByUserName 根据name获取user
	GetUserByUserName(name string) *model.User
	//GetUserByUserId 根据id获取user
	GetUserByUserId(id int64) *model.User
	//AddUser 添加用户
	AddUser(user model.User) bool
	//GetUserRespondById 根据id获取User对象
	GetUserRespondById(id int64) *model.User
}
type UserServiceImpl struct {
	C *gin.Context
	UserService
}
