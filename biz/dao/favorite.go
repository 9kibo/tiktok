package dao

import (
	"tiktok/biz/model"
	"tiktok/pkg/errno"
)

// GetFavoriteVideoIdS 根据User_id获取用户点赞的所有Video_Id
func GetFavoriteVideoIdS(UserId int64) ([]int64, error) {
	var videoIDs []int64

	err := Db.Table("video_favor").Where("user_id = ?", UserId).Pluck("video_id", &videoIDs).Error
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
		return errno.FavoriteRelationAlreadyExistErr
	}
	return nil

}
