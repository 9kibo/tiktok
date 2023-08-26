package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tiktok/biz/model"
)

func TestAddFollow(t *testing.T) {
	assert.NoError(t, AddFollow(&model.Follow{
		FolloweeId: 1,
		FollowerId: 2,
	}))
}
func TestDeleteFollow(t *testing.T) {
	assert.NoError(t, AddFollow(&model.Follow{
		FolloweeId: 11,
		FollowerId: 22,
	}))
	assert.NoError(t, DeleteFollow(&model.Follow{
		FolloweeId: 11,
		FollowerId: 22,
	}))
}
func TestExistsFollow(t *testing.T) {
	assert.NoError(t, AddFollow(&model.Follow{
		FolloweeId: 111,
		FollowerId: 222,
	}))
	exists, err := ExistsFollow(&model.Follow{
		FolloweeId: 111,
		FollowerId: 222,
	})
	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, DeleteFollow(&model.Follow{
		FolloweeId: 111,
		FollowerId: 222,
	}))

	exists, err = ExistsFollow(&model.Follow{
		FolloweeId: 111,
		FollowerId: 222,
	})
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestGetFollowerList(t *testing.T) {
	var followeeId int64 = 111
	var followerId int64 = 222

	assert.NoError(t, AddFollow(&model.Follow{
		FolloweeId: followeeId,
		FollowerId: followerId,
	}))
	followings, err := GetFollowerList(followeeId)
	assert.NoError(t, err)
	assert.True(t, followings[0].FollowerId == followerId)
	assert.NoError(t, DeleteFollow(&model.Follow{
		FolloweeId: followeeId,
		FollowerId: followerId,
	}))
}
func TestGetFolloweeList(t *testing.T) {
	var followeeId int64 = 111
	var followerId int64 = 222

	assert.NoError(t, AddFollow(&model.Follow{
		FolloweeId: followeeId,
		FollowerId: followerId,
	}))
	followings, err := GetFollowList(followerId)
	assert.NoError(t, err)
	assert.True(t, followings[0].FolloweeId == followeeId)
	assert.NoError(t, DeleteFollow(&model.Follow{
		FolloweeId: followeeId,
		FollowerId: followerId,
	}))
}

func TestGetFollowingCount(t *testing.T) {
	_, err := GetFollowingCount(222)
	assert.NoError(t, err)
}
func TestGetFollowerCount(t *testing.T) {
	_, err := GetFollowerCount(222)
	assert.NoError(t, err)
}
