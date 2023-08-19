package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/kafka"
	"tiktok/biz/middleware/logmw"
	tredis "tiktok/biz/middleware/redis"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
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
	FavouriteAction(userId int64, videoId int64, actionType int32) error
	// GetFavouriteList 获取当前用户的所有点赞视频
	GetFavouriteList(userId int64, curId int64) ([]*model.Video, error)
}

type FavoriteImpl struct {
	C *gin.Context
}

func NewFavorite(c *gin.Context) *FavoriteImpl {
	return &FavoriteImpl{
		C: c,
	}
}

func (F *FavoriteImpl) FavouriteAction(userId int64, videoId int64, actionType int32) error {
	userIdStr := strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	var Rctx = context.Background()
	rdb, err := tredis.GetRedis(8)
	//建立redis连接
	defer rdb.Close()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis连接错误")
		return err
	}
	//写入消息失败的回调函数
	addBack := func(k string, v string) {
		logmw.LogWithRequest("Favorite", F.C).Debug("kafka写入失败")
		rdb.SRem(Rctx, userIdStr, videoId)
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":Kafka写入失败"))
	}
	delBack := func(k string, v string) {
		logmw.LogWithRequest("Favorite", F.C).Debug("kafka写入失败")
		rdb.SAdd(Rctx, userIdStr, videoId)
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":Kafka写入失败"))
	}
	//判断缓存中key是否存在
	exists, err := rdb.Exists(Rctx, userIdStr).Result()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis连接错误")
		return err
	}
	if exists < 1 {
		//缓存中不存在,更新缓存
		err = LoadFavoriteToRides(userId, rdb, Rctx)
		if err != nil {
			F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
			logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis更新错误")
			return err
		}
	} else {
		//缓存中存在，更新缓存时间
		if ttl, _ := rdb.TTL(Rctx, userIdStr).Result(); ttl < constant.Favorite_UserId_DefaultTime/3 {
			rdb.Expire(Rctx, userIdStr, constant.Favorite_UserId_DefaultTime/3)
		}
	}
	var n int64
	if actionType == 1 {
		//点赞
		n, err = rdb.SAdd(Rctx, userIdStr, videoIdStr).Result()
		if err != nil {
			F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
			logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis写入失败")
			return err
		}
		if n == 0 {
			logmw.LogWithRequest("Favorite", F.C).Debug("重复点赞")
			return nil
		}
		//更新redis成功，向消息队列发送key:UserId  value:videoId+time.now().unix(),为方便直接使用字符串拼接，追求性能可以考虑使用结构体将结构体序列化
		value := videoIdStr + " " + strconv.FormatInt(time.Now().Unix(), 10)
		kafka.FavoriteMq.WriteMsg(userIdStr, value, addBack)
		return nil
	} else {
		//取消点赞
		n, err = rdb.SRem(Rctx, userIdStr, videoId).Result()
		if err != nil {
			F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
			logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis删除失败")
			return err
		}
		if n == 0 {
			logmw.LogWithRequest("Favorite", F.C).Debug("重复取消赞")
			return nil
		}
		//更新完redis后向消息队列推送,为防止消息乱序，推送到同一个topic，根据是否有createAt时间戳来判断是删除还是增加
		//后期优化可以考虑增加topic分片，增加消费者，使用userId作为hash对象仍然可以保证同一用户的消息顺序性
		kafka.FavoriteMq.WriteMsg(userIdStr, strconv.FormatInt(videoId, 10), delBack)
	}
	return nil
}

// LoadFavoriteToRides 从将点赞信息加载到Redis中
func LoadFavoriteToRides(UserId int64, rdb *redis.Client, Rctx context.Context) error {
	VideoIds, err := dao.GetFavoriteVideoIdS(UserId)
	if err != nil {
		return err
	}
	UserIdStr := strconv.FormatInt(UserId, 10)
	//添加占位，防止点赞全部取消后key被删除
	err = rdb.SAdd(Rctx, UserIdStr, -1).Err()
	err = rdb.SAdd(Rctx, UserIdStr, VideoIds).Err()
	if err != nil {
		return err
	}
	//设置过期时间，默认3天
	err = rdb.Expire(Rctx, UserIdStr, constant.Favorite_UserId_DefaultTime).Err()
	if err != nil {
		return err
	}
	return nil
}
