package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserId  uint `gorm:"not null;index:idx_favorite" json:"user_id"`
	VideoId uint `gorm:"not null;index:idx_favorite" json:"video_id"`
	Cancel  bool `gorm:"" json:"cancel"`
}
