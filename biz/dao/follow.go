package dao

import (
	"tiktok/biz/model"
)

// AddFollow 添加关注
func AddFollow(follow *model.Follow) error {
	return ofUpdate1(Db.Create(follow))
}

// DeleteFollow 删除关注
func DeleteFollow(follow *model.Follow) error {
	return ofUpdate1(Db.Where("followee_id=? and follower_id=?",
		follow.FolloweeId, follow.FollowerId).Delete(&model.Follow{}))
}

// ExistsFollow 有数据库错误才返回错误
func ExistsFollow(follow *model.Follow) (bool, error) {
	return ofExists(Db.Where("followee_id=? and follower_id=?",
		follow.FolloweeId, follow.FollowerId).Take(&model.Follow{}).Error)
}

// GetFollowerList 获取粉丝
func GetFollowerList(userId int64) ([]*model.Follow, error) {
	var follows []*model.Follow
	err := Db.Select("follower_id").Find(&follows, "followee_id=?", userId).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

// GetFollowList 获取关注的人
func GetFollowList(userId int64) ([]*model.Follow, error) {
	var follows []*model.Follow
	tx := Db.Select("followee_id").Find(&follows, "follower_id=?", userId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return follows, nil
}

func GetFollowingCount(userId int64) (int64, error) {
	var count int64
	tx := Db.Model(&model.Follow{}).Where("follower_id=?", userId).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}

func GetFollowerCount(userId int64) (int64, error) {
	var count int64
	tx := Db.Model(&model.Follow{}).Where("followee_id=?", userId).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return count, nil
}
