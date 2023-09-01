package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/kafka"
	tredis "tiktok/biz/middleware/redis"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	"time"
)

type FavoriteService interface {
	//IsFavorite 根据id返回是否点赞了该视频
	IsFavorite(videoId int64, userId int64) (bool, error)
	//FavouriteCount 根据当前视频id获取当前视频点赞数量。
	FavouriteCount(videoId int64) (int64, error)
	//TotalFavourite 根据userId获取这个用户总共被点赞数量
	TotalFavourite(userId int64) (int64, error)
	//FavouriteVideoCount 根据userId获取这个用户点赞视频数量
	FavouriteVideoCount(userId int64) (int64, error)

	//FavouriteAction 当前操作行为，1点赞，2取消点赞。
	FavouriteAction(userId int64, videoId int64, actionType int32)
	// GetFavouriteList 获取当前用户的所有点赞视频
	GetFavouriteList(userId int64, curId int64) []*model.Video
}

type FavoriteImpl struct {
	FavoriteService
	C *gin.Context
}

func NewFavorite(c *gin.Context) *FavoriteImpl {
	return &FavoriteImpl{
		C: c,
	}
}

// FavouriteAction 点赞信息并非重要信息，可以允许一段时间的不一致，
// 为了保证接口的响应速度，redis和数据库的更新均为异步进行，将消息成功写入消息队列即可认为操作成功，后续错误打印日志
func (F *FavoriteImpl) FavouriteAction(userId int64, videoId int64, actionType int32) {
	//判断视频是否存在
	VideoS := &VideoServiceImpl{ctx: F.C}
	if _, err := VideoS.GetVideoById(videoId, userId); err != nil {
		utils.LogBizErr(F.C, errno.VideoIsNotExistErr, http.StatusOK, "视频不存在")
	}
	rdb, err := tredis.FavR.GetFavRedis()
	if err != nil {
		utils.LogDB(F.C, errno.Service)
	}
	//将消息写入队列
	msg := strings.Builder{}
	msg.WriteString(strconv.FormatInt(videoId, 10))
	now := time.Now().Unix()
	//点赞操作增加创建时间
	if actionType == 1 {
		msg.WriteString(" ")
		msg.WriteString(strconv.FormatInt(now, 10))
	}
	kafka.FavoriteMq.WriteMsg(strconv.FormatInt(userId, 10), msg.String(), func(string, string) {
		utils.LogBizErr(F.C, errno.FavoriteActionErr, http.StatusOK, "kafka写入失败")
	})
	//异步更新缓存
	go tredis.UpdateFavRedis(rdb, userId, videoId, actionType, now, F.C.GetString(constant.RequestId))
}

// GetFavouriteList 获取点赞视频列表
func (F *FavoriteImpl) GetFavouriteList(userId int64, curId int64) []*model.Video {
	VideoServer := &VideoServiceImpl{
		ctx: F.C,
	}
	rdb, err := tredis.FavR.GetFavRedis()
	if err != nil {
		utils.LogDB(F.C, errno.Service)
	}
	//判断缓存
	err = tredis.LoadIfNotExists(userId, rdb, tredis.LoadFavoriteToRides)
	if err != nil {
		utils.LogBizErr(F.C, errno.Update, http.StatusOK, "更新redis出错")
	}
	//获取所有点赞列表
	var VideoIds []int64
	vals, _ := rdb.ZRevRange(context.Background(), strconv.FormatInt(userId, 10), 0, -1).Result()
	for _, val := range vals {
		id, _ := strconv.ParseInt(val, 10, 64)
		if id == -1 {
			//占位符
			continue
		}
		VideoIds = append(VideoIds, id)
	}
	VideoList := make([]*model.Video, len(VideoIds))
	var wg sync.WaitGroup
	for i, VideoId := range VideoIds {
		wg.Add(1)
		go func(a int, id int64) {
			Video, err := VideoServer.GetVideoById(id, curId)
			if err != nil {
				utils.LogWithRequestId(F.C, "Favorite", err).WithField("videoId:", id).Error("视频获取出错")
				wg.Done()
				return
			}
			VideoList[a] = Video
			wg.Done()
		}(i, VideoId)
	}
	wg.Wait()
	return VideoList
}

// IsFavorite 判断是否点赞该视频
func (F *FavoriteImpl) IsFavorite(videoId int64, userId int64) (bool, error) {
	var Rctx = context.Background()
	rdb, err := tredis.FavR.GetFavRedis()
	if err != nil {
		return false, err
	}
	exists, err := rdb.Exists(Rctx, strconv.FormatInt(userId, 10)).Result()
	if err != nil {
		return false, err
	}
	//缓存中不存在去数据库中查
	if exists < 1 {
		return dao.ExistsFav(userId, videoId)
	}
	val, _ := rdb.ZScore(Rctx, strconv.FormatInt(userId, 10), strconv.FormatInt(videoId, 10)).Result()
	if val == float64(0) {
		return false, nil
	}
	return true, nil
}

// FavouriteVideoCount 根据userId获取这个用户点赞视频数量
func (F *FavoriteImpl) FavouriteVideoCount(userId int64) (int64, error) {
	rdb, err := tredis.FavR.GetFavRedis()
	if err != nil {
		return 0, err
	}
	if err = tredis.LoadIfNotExists(userId, rdb, tredis.LoadFavoriteToRides); err != nil {
		return 0, err
	}
	n, err := rdb.ZCard(context.Background(), strconv.FormatInt(userId, 10)).Result()
	n = n - 1 //减去占位符
	return n, err
}

// FavouriteCount 根据当前视频id获取当前视频点赞数量。
func (F *FavoriteImpl) FavouriteCount(videoId int64) (int64, error) {
	//缓存中 videoId->count
	videoIdStr := strconv.FormatInt(videoId, 10)
	var Rctx = context.Background()
	rdb, err := tredis.FavR.GetFavRedis()
	if err != nil {
		return 0, err
	}
	err = tredis.LoadIfNotExists(videoId, rdb, tredis.LoadVideoFavCount)
	if err != nil {
		return 0, err
	}
	n, err := rdb.Get(Rctx, videoIdStr).Int64()
	if err != nil {
		//兜底机制
		n, err = dao.GetVideoFavorCount(videoId)
		if err != nil {
			return n, err
		}
		return n, nil
	}
	return n, nil
}
