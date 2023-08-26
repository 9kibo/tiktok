package dao

import (
	"fmt"
	"tiktok/biz/model"
)

func AddUser(user *model.User) (int64, error) {
	tx := Db.Create(user)
	return user.Id, ofUpdate1(tx)
}

func ExistsUserByUsername(username string) (bool, error) {
	return ofExists(Db.Where("username=?", username).Take(&model.User{}).Error)
}
func ExistsUserById(userId int64) (bool, error) {
	return ofExists(Db.Take(&model.User{
		Id: userId,
	}).Error)
}

func MustGetUserByUsernamePassword(username, password string) (*model.User, error) {
	user := model.User{}
	tx := Db.Where("username=? and password=?", username, password).Take(&user)
	if err := ofGet(tx.Error); err != nil {
		return nil, err
	}
	return &user, nil
}

func MustGetUserById(userId int64) (*model.User, error) {
	user := model.User{
		Id: userId,
	}
	tx := Db.Take(&user)
	if err := ofGet(tx.Error); err != nil {
		return nil, err
	}
	return &user, nil
}

func MustGetUsersByIds(userIds []int64) ([]*model.User, error) {
	users := make([]*model.User, 0, len(userIds))
	r := Db.Where("id in (?)", userIds).Order("username").Find(&users)
	if r.Error != nil {
		return nil, r.Error
	}
	if len(users) != len(userIds) {
		m := make(map[int64]struct{}, len(userIds))
		for _, user := range users {
			m[user.Id] = struct{}{}
		}
		notFindUserIds := make([]int64, 0, len(userIds)-len(users))
		for _, userId := range userIds {
			if _, ok := m[userId]; !ok {
				notFindUserIds = append(notFindUserIds, userId)
			}
		}
		return nil, fmt.Errorf("can't find in db, userIds=[%v]", notFindUserIds)
	}
	return users, nil
}
