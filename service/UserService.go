package service

import (
	"tiktok/model"
)

type UserRespond struct {
	Id              int64  `json:"id"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	Signature       string `json:"signature,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
}
type UserService interface {
	//GetUserByUserName 根据name获取user
	GetUserByUserName(name string) model.User
	//GetUserByUserId 根据id获取user
	GetUserByUserId(id int64) model.User
	//AddUser 添加用户
	AddUser(NewUser model.User) bool
	//GetUserRespondById 根据id获取User对象
	GetUserRespondById(id int64) UserRespond
}
