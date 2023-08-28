package dao

import (
	"fmt"
	"testing"
	"tiktok/biz/model"
	"time"
)

func TestFavorite(t *testing.T) {

	// 准备测试数据
	var testData = []model.VideoFavor{
		{UserId: 1, VideoId: 1},
		{UserId: 1, VideoId: 2},
		{UserId: 1, VideoId: 3},
		{UserId: 1, VideoId: 4},
		{UserId: 1, VideoId: 5},
		{UserId: 1, VideoId: 6},
		{UserId: 1, VideoId: 7},
		{UserId: 1, VideoId: 8},
		{UserId: 1, VideoId: 9},
		{UserId: 1, VideoId: 10},
	}

	// 添加点赞
	for _, data := range testData {
		if err := AddFavorite(data.UserId, data.VideoId, time.Now().Unix()); err != nil {
			println("_________________")
			fmt.Printf("%s\n", err)
			println("_________________")
		}
	}

	// 查询点赞视频
	user1Favs, _ := GetUserFavorCount(1)
	_, err := GetUserFavorCount(2)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// 验证
	if user1Favs != 10 {
		t.Errorf("expect 10 favorite for user 1, but got %d", user1Favs)
	}

	// 删除点赞
	for _, data := range testData {
		if err := DelFavorite(data.UserId, data.VideoId); err != nil {
			fmt.Printf("%s\n", err)
		}
	}

	// 查询点赞视频
	user1Favs, _ = GetUserFavorCount(1)

	// 验证删除
	if user1Favs != 0 {
		t.Errorf("expect 0 favorite after delete, but got %d", user1Favs)
	}
}
