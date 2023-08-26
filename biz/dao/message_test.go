package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tiktok/biz/model"
)

func TestAddMessage(t *testing.T) {
	err := AddMessage(&model.Message{
		FromUserId: 1,
		ToUserId:   2,
		Content:    "sb",
	})
	assert.NoError(t, err)

	err = AddMessage(&model.Message{
		FromUserId: 2,
		ToUserId:   1,
		Content:    "sb also",
	})
	assert.NoError(t, err)
}
func TestMustGetMessageListByLastTime(t *testing.T) {
	messageList, err := GetMessageRecordByLastTime(1, 2, 1692861028)
	assert.NoError(t, err)
	assert.True(t, len(messageList) == 2)
}
