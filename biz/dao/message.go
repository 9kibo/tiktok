package dao

import (
	"fmt"
	"sync"
	"tiktok/biz/model"
)

func AddMessage(message *model.Message) error {
	return ofUpdate1(Db.Create(message))
}
func GetMessageRecordByLastTime(userId, toUserId, lastTime int64) ([]*model.Message, error) {
	messageList := make([]*model.Message, 0, 2)
	tx := Db.Where("created_at > ? and ((from_user_id = ? and to_user_id = ?) or (from_user_id = ? and to_user_id = ?))",
		lastTime, userId, toUserId, toUserId, userId).Order("created_at").Find(&messageList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return messageList, nil
}
func GetFriendIds(userId int64) ([]int64, error) {
	messageList := make([]*model.Message, 0, 2)
	tx := Db.Distinct("from_user_id", "to_user_id").Where("from_user_id = ? or to_user_id = ?", userId, userId).Find(&messageList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	friendIds := make([]int64, 0, len(messageList))
	mmp := make(map[int64]struct{}, len(messageList))
	for _, message := range messageList {
		if message.FromUserId == userId {
			mmp[message.ToUserId] = struct{}{}
			friendIds = append(friendIds, message.ToUserId)
		} else if _, exists := mmp[message.FromUserId]; !exists && message.ToUserId == userId {
			mmp[message.FromUserId] = struct{}{}
			friendIds = append(friendIds, message.FromUserId)
		}
	}
	return friendIds, nil
}

// GetMessageLatestByToUserIds 一对一地并发地查
// 使用context和Db.WithContext(ctx)可以控制失败一个就取消其他任务, 但是暂时先不支持, 容错?
func GetMessageLatestByToUserIds(userId int64, toUserIds []int64) ([]*model.Message, error) {
	var err error = nil
	messageList := make([]*model.Message, len(toUserIds))
	wg := sync.WaitGroup{}
	//ctx, cancelFunc := context.WithCancel(context.Background())
	wg.Add(len(toUserIds))
	for i, toUserId := range toUserIds {
		go func(i int, toUserId int64) {
			defer wg.Done()
			message := model.Message{}
			tx := Db.Select("content").Where(
				"(from_user_id = ? and to_user_id =?) or (to_user_id = ? and from_user_id =?)",
				userId, toUserId, toUserId, userId).Order(
				"created_at desc").Take(&message)
			if tx.Error != nil {
				//cancelFunc()
				err = fmt.Errorf("when get user[%d] and toUser[%d] chat record, err=%s", userId, toUserId, tx.Error)
				return
			}
			messageList[i] = &message
		}(i, toUserId)
	}
	wg.Wait()
	return messageList, err
}
