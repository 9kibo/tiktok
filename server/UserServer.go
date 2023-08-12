package server

import (
	"tiktok/model"
)

type UserRespond struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
}
type UserServer interface {
	//GetUserByUserName 根据name获取user
	GetUserByUserName(name string) model.User
	//GetUserByUserId 根据id获取user
	GetUserByUserId(id int64) model.User
	//AddUser 添加用户
	AddUser(NewUser model.User) bool
	//GetUserRespondById 根据id获取User对象
	GetUserRespondById(id int64) UserRespond
}
