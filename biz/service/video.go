package service

import (
	"mime/multipart"
	"tiktok/biz/model"
	"time"
)

type VideoService interface {
	//Feed 传入时间戳,当前用户Id，返回视频切片 和 返回的视频切片中的最早时间
	Feed(lastTime time.Time, userId int64) ([]*model.Video, time.Time, error)
	//GetVideoById 根据视频id和用户id获取video
	GetVideoById(videoId int64, userId int64) (*model.Video, error)
	//Publish 上传视频
	Publish(file multipart.File, fileHeader *multipart.FileHeader, userId int64, title string) error
	//GetVideoListById 当前用户 (userId) 获取目标用户(targetId)发布的视频
	GetVideoListById(targetId int64, userId int64)
}
type VideoServiceImpl struct {
}
