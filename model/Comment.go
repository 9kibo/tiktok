package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId      uint   `gorm:"not null"`
	VideoId     uint   `gorm:"not null"`
	CommentText string `gorm:"not null;type:longtext"`
	cancel      bool   `gorm:"not null"`
}
