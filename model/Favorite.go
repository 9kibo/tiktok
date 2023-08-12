package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	UserId  uint `gorm:"not null;index:idx_favorite"`
	VideoId uint `gorm:"not null;index:idx_favorite"`
	Cancel  bool `gorm:"not null"`
}
