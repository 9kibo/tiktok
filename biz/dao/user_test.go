package dao

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"tiktok/biz/model"
	"tiktok/pkg/errno"
)

func TestAddUser(t *testing.T) {
	userId, err := AddUser(&model.User{
		Username: "dao-test1",
		Password: "123",
	})
	assert.NoError(t, err)
	assert.True(t, userId > 0)
}
func TestExistsUserByUsername(t *testing.T) {
	exists, err := ExistsUserByUsername("dao-test1")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = ExistsUserByUsername("XXXXXXX")
	assert.NoError(t, err)
	assert.False(t, exists)
}
func TestGetUserIdByUsernamePassword(t *testing.T) {
	user, err := MustGetUserByUsernamePassword("dao-test1", "123")
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = MustGetUserByUsernamePassword("dao-test12", "123")
	assert.Error(t, err)
	assert.Equal(t, errno.NotExists, err)
	assert.Nil(t, user)
}
func TestGetUserById(t *testing.T) {
	user, err := MustGetUserById(4)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = MustGetUserById(1231231231233123123)
	assert.Error(t, err)
	assert.Equal(t, errno.NotExists, err)
	assert.Nil(t, user)
}
func TestGetUsersByIds(t *testing.T) {
	users, err := MustGetUsersByIds([]int64{3, 4})
	assert.NoError(t, err)
	assert.Equal(t, len(users), 2)

	users, err = MustGetUsersByIds([]int64{3, 4, 5, 6})
	assert.Error(t, err)
	assert.Equal(t, fmt.Sprintf("can't find in db, userIds=[%v]", []int{5, 6}), err.Error())
	assert.Nil(t, users)
}
