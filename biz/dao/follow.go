package dao

import (
	"errors"
	"gorm.io/gorm"
	"tiktok/biz/model"
	"tiktok/pkg/errno"
)

// AddFollow 添加关注
func AddFollow(f *model.Follow) error {
	return Db.Create(f).Error
}

// DeleteFollow 删除关注
func DeleteFollow(f *model.Follow) error {
	tx := Db.Where("followee_id=? and follower_id=?",
		f.FolloweeId, f.FollowerId).Delete(&model.Follow{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errno.FollowRelationAlreadyExistErr
	}
	return nil
}

// ExistsFollow 有数据库错误才返回错误
func ExistsFollow(f *model.Follow) (bool, error) {
	err := Db.Where("followee_id=? and follower_id=?",
		f.FolloweeId, f.FollowerId).Take(&model.Follow{}).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}

// GetFollowerList 获取粉丝
func GetFollowerList(f *model.Follow) ([]*model.Follow, error) {
	var follows []*model.Follow
	err := Db.Model(f).Select("follower_id").Find(&follows, "followee_id=?", f.FolloweeId).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

// GetFolloweeList 获取关注的人
func GetFolloweeList(f *model.Follow) ([]*model.Follow, error) {
	var follows []*model.Follow
	err := Db.Model(f).Select("followee_id").Find(&follows, "follower_id=?", f.FollowerId).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}
