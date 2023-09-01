package dao

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tiktok/biz/model"
	"time"
)

func TestVideoDAO_Create(t *testing.T) {
	tmpVideo := model.Video{
		CreatedAt:     time.Now().Unix(),
		UpdatedAt:     0,
		DeletedAt:     0,
		AuthorId:      11,
		Title:         "title",
		PlayUrl:       "videoUrl",
		CoverUrl:      "coverUrl",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Author:        nil,
	}

	assert.NoError(t, NewVideo().Create(&tmpVideo))
	t.Logf("videoId=%#v", tmpVideo)
}

func TestVideoDAO_Delete(t *testing.T) {
	tmpVideo := model.Video{
		Id: 1,
	}
	// 如果删除不存在的数据，也会返回Success
	assert.NoError(t, NewVideo().Delete(&tmpVideo))
	t.Logf("Successful Deleted")
}

func TestVideoDAO_GetVideoById(t *testing.T) {
	tmpVideo := model.Video{
		Id: 2,
	}
	nTmpVideo, err := NewVideo().GetVideoById(&tmpVideo)
	assert.NoError(t, err)
	assert.True(t, nTmpVideo.Id == 2)
	t.Logf("videoId=%#v", nTmpVideo)
}

func TestVideoDAO_UpdateVideo(t *testing.T) {
	tmpVideo := model.Video{
		Id: 2,
	}
	// 查询
	nTmpVideo, err := NewVideo().GetVideoById(&tmpVideo)
	assert.NoError(t, err)
	t.Logf("videoId=%#v", nTmpVideo)
	// 修改
	nTmpVideo.Title = "Hello World"
	assert.NoError(t, NewVideo().UpdateVideo(nTmpVideo, nTmpVideo))
	t.Logf("videoId=%#v", nTmpVideo)
}

func TestVideoDAO_GetVideoListById(t *testing.T) {
	idList := []int64{2, 3}
	tmpVideoList, err := NewVideo().GetVideoListById(idList)
	assert.NoError(t, err)
	t.Logf("VideoList=%#v", tmpVideoList)
	for i := range tmpVideoList {
		t.Logf("VideoList[i]=%#v", tmpVideoList[i])
	}
}

func TestVideoDAO_GetVideoListByLastTime(t *testing.T) {
	limit := 30
	tmpVideo, err := NewVideo().GetVideoListByLastTime(time.Unix(1692867192, 0), limit)
	assert.NoError(t, err)
	assert.False(t, len(tmpVideo) == 0)
	t.Logf("VedioList=%#v", tmpVideo)
	for i := range tmpVideo {
		t.Logf("VideoList[i]=%#v", tmpVideo[i])
	}
}
