package dao

import "tiktok/biz/model"

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
func GetMessageLatestByToUserIds(userId int64, toUserIds []int64) ([]*model.Message, error) {
	messageList := make([]*model.Message, 0, 2)
	tx := Db.Where("(from_user_id = ? and to_user_id in (?)) or (from_user_id = ? and to_user_id in (?))",
		userId, toUserIds, toUserIds, userId).Order("created_at").Find(&messageList)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return messageList, nil
}
