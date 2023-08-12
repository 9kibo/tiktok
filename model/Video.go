package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	VideoAuthorId uint   `gorm:"index;not null"`
	VideoUrl      string `gorm:"" json:"video_url"`
	VideoCover    string `gorm:"" json:"video_cover"`
	Tittle        string `gorm:"" json:"tittle"`
}
