package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(20);not null" json:"username"`
	Password string `gorm:"type:varchar(500);not null" json:"password"`
}
