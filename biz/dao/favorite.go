package dao

import (
	"tiktok/biz/model"
	"tiktok/pkg/errno"
)

// GetFavoriteVideoIdS 根据User_id获取用户点赞的所有Video_Id和CreateAt
func GetFavoriteVideoIdS(UserId int64) ([]model.VideoFavor, error) {
	var videoIDs []model.VideoFavor

	err := Db.Table("video_favor").Where("user_id = ?", UserId).Select("video_id, created_at").Scan(&videoIDs).Error
	if err != nil {
		if err.Error() == "record not found" {
			return videoIDs, nil
		}
		return nil, err
	}
	return videoIDs, nil
}

// AddFavorite 添加点赞信息
func AddFavorite(userId int64, videoId int64, createAt int64) error {
	var favorite model.VideoFavor
	favorite = model.VideoFavor{
		UserId:    userId,
		VideoId:   videoId,
		CreatedAt: createAt,
	}
	tx := Db.Table("video_favor").Create(&favorite)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errno.FavoriteRelationAlreadyExistErr
	}
	return nil

}

// DelFavorite 删除点赞信息
func DelFavorite(userId int64, videoId int64) error {
	var favorite model.VideoFavor
	tx := Db.Table("video_favor").Where("user_id = ? AND video_id = ?", userId, videoId).Delete(&favorite)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errno.FavoriteRelationNotExistErr
	}
	return nil

}

// 判断是否点赞
func ExistsFav(userId int64, videoId int64) (bool, error) {
	var fav model.VideoFavor
	err := Db.Where("user_id = ? AND video_id = ?", userId, videoId).First(&fav).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 获取用户点赞数
func GetUserFavorCount(userId int64) (int64, error) {
	var count int64
	if err := Db.Model(&model.VideoFavor{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		if err.Error() == "record not found" {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

// 获取视频点赞数
func GetVideoFavorCount(videoId int64) (int64, error) {
	var count int64
	if err := Db.Model(&model.VideoFavor{}).Where("video_id = ?", videoId).Count(&count).Error; err != nil {
		if err.Error() == "record not found" {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}
