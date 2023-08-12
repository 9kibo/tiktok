package model

import (
	"gorm.io/gorm"
)

type Follow struct {
	gorm.Model
	UserId   uint `gorm:"not null;index:idx_follower"` //用户id
	FollowId uint `gorm:"not null;index:idx_follower"` //UserId关注的用户的Id
	Cancel   bool `gorm:"not null"`
}
