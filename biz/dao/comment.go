package dao

import (
	"github.com/pkg/errors"
	"tiktok/biz/model"
)

// 添加评论
func AddComm(comm *model.Comment) error {
	return Db.Create(comm).Error
}

// 删除评论
func DelComm(commId, userId int64) error {
	// 查询评论
	var comm model.Comment
	Db.First(&comm, commId)

	// 验证用户ID
	if comm.UserId != userId {
		return errors.New("用户无权删除评论")
	}

	return Db.Delete(&comm).Error
}

// 获取评论数量
func GetCommCount(videoId int64) (int64, error) {
	var count int64
	err := Db.Model(&model.Comment{}).Where("video_id = ?", videoId).Count(&count).Error
	return count, err
}

// 获取评论列表
func GetCommList(videoId int64) ([]*model.Comment, error) {
	var List []*model.Comment
	err := Db.Where("video_id = ?", videoId).Find(&List).Error
	return List, err
}

// 根据ID获取评论
func Comm(commId int64) (*model.Comment, error) {
	comm := model.Comment{}
	err := Db.Find(&comm, commId).Error
	return &comm, err
}
