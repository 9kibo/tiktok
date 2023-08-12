package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	VideoAuthorId uint   `gorm:"not null;index;"`
	VideoUrl      string `gorm:"not null;"`
	VideoCover    string `gorm:"not null;"`
	Tittle        string `gorm:"not null;type:longtext"`
}
