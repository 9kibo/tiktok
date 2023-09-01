package dao

import (
	"errors"
	"gorm.io/gorm"
	"sync"
	"tiktok/biz/model"
	"time"
)

type VideoDAO struct{}

var (
	videoDAOHandler *VideoDAO
	videoDAOOnce    sync.Once
)

// NewVideo 单例实现
func NewVideo() *VideoDAO {
	videoDAOOnce.Do(func() {
		videoDAOHandler = &VideoDAO{}
	})
	return videoDAOHandler
}

// Create 新增视频数据
func (dao *VideoDAO) Create(v *model.Video) error {
	return Db.Model(&model.Video{}).Create(v).Error
}

func (dao *VideoDAO) Delete(v *model.Video) error {
	return Db.Model(&model.Video{}).Delete(v).Error
}

// GetVideoById 根据视频ID获取视频信息
// 不包含User和点赞信息，可以在Service层调用其他的DAO方法
func (dao *VideoDAO) GetVideoById(v *model.Video) (*model.Video, error) {
	var targetVideo model.Video
	// 获取Video对象
	err := Db.Model(v).First(&targetVideo, v.Id).Error
	//// 根据视频的author_id 获取作者信息
	//targetAuthor := model.User{}
	//err = Db.Model(&model.User{}).
	//	Omit("password", "created_at", "updated_at").
	//	Where("id = ?", targetVideo.AuthorId).First(&targetAuthor).Error
	//targetVideo.Author = &targetAuthor

	if err != nil {
		return nil, err
	}
	return &targetVideo, nil
}

// GetVideoListById 批量获取视频信息
func (dao *VideoDAO) GetVideoListById(videoIds []int64) ([]*model.Video, error) {
	var targetVideoList []*model.Video
	err := Db.Model(&model.Video{}).Where("Id IN ?", videoIds).Find(&targetVideoList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到记录，返回空切片
			return []*model.Video{}, nil
		}
		return nil, err
	}
	return targetVideoList, nil
}

// UpdateVideo 更新视频信息
func (dao *VideoDAO) UpdateVideo(v *model.Video, n *model.Video) error {
	err := Db.Model(v).Updates(n).Error
	if err != nil {
		return err
	}
	return nil
}

// GetVideoListByLastTime 根据数量要求返回最新创建的视频列表
func (dao *VideoDAO) GetVideoListByLastTime(lastTime time.Time, limit int) ([]*model.Video, error) {
	var targetVideoList []*model.Video
	err := Db.Model(&model.Video{}).
		Where("created_at < ?", lastTime.Unix()).
		Order("created_at desc").
		Limit(limit).
		Find(&targetVideoList).Error
	if err != nil {
		return nil, err
	}
	return targetVideoList, nil
}

// GetVideoListByAuthor 根据作者Id查询视频
func (dao *VideoDAO) GetVideoListByAuthor(authorId int64) ([]*model.Video, error) {
	var targetVideoList []*model.Video
	err := Db.Model(&model.Video{}).
		Where("author_id = ?", authorId).
		Order("created_at desc").
		Find(&targetVideoList).Error
	if err != nil {
		return nil, err
	}
	return targetVideoList, nil
}
