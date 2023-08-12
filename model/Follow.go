package model

import (
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	UserId     uint `gorm:"not null;index:idx_follower" json:"user_id"`  //用户id
	FollowerId uint `gorm:"not null;index:idx_follower" json:"video_id"` //关注的用户
	Cancel     bool `gorm:"" json:"cancel"`
}
