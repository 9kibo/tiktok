package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName   string `gorm:"not null;index:idx_name_pwd"`
	Password   string `gorm:"not null;index:idx_name_pwd"`
	Avatar     string `gorm:"default:https://th.bing.com/th/id/OIP.FlsXXU-wKCf-Us4zDBfvzwHaHa?pid=ImgDet&rs=1"` //头像
	Background string `gorm:"default:https://th.bing.com/th/id/OIP.3C4Yst9hMZpZFOh2kCzNFwAAAA?pid=ImgDet&rs=1"` //背景图片
	signature  string `gorm:"default:无个人简介"`                                                                    //个人简介
}
