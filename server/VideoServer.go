package server

import (
	"mime/multipart"
	"time"
)

type VideoRespond struct {
	Id            int64       `json:"id"`
	Author        UserRespond `json:"author"`
	PlayUrl       string      `json:"play_url"`
	CoverUrl      string      `json:"cover_url"`
	FavoriteCount int64       `json:"favorite_count"`
	CommentCount  int64       `json:"comment_count"`
	IsFavorite    bool        `json:"is_favorite"`
	Title         string      `json:"title"`
}

type VideoServer interface {
	//Feed 传入时间戳,当前用户Id，返回视频切片 和 返回的视频切片中的最早时间
	Feed(lastTime time.Time, userId int64) ([]VideoRespond, time.Time, error)
	//GetVideoById 根据视频id和用户id获取video
	GetVideoById(videoId int64, userId int64) (VideoRespond, error)
	//Publish 上传视频
	Publish(file multipart.File, fileHeader *multipart.FileHeader, userId int64, title string) error
	//GetVideoListById 当前用户 (userId) 获取目标用户(targetId)发布的视频
	GetVideoListById(targetId int64, userId int64)
}
