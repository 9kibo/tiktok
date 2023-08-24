package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
	"sync"
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
	FavoriteService
	C *gin.Context
}

func NewFavorite(c *gin.Context) *FavoriteImpl {
	return &FavoriteImpl{
		C: c,
	}
}

func (F *FavoriteImpl) FavouriteAction(userId int64, videoId int64, actionType int32) error {
	//判断视频是否存在
	VideoS := &VideoServiceImpl{C: F.C}
	if _, err := VideoS.GetVideoById(videoId, userId); err != nil {
		F.C.AbortWithStatusJSON(http.StatusBadRequest, errno.NewErrno(errno.VideoIsNotExistErrCode, "视频不存在"))
	}
	userIdStr := strconv.FormatInt(userId, 10)
	videoIdStr := strconv.FormatInt(videoId, 10)
	var Rctx = context.Background()
	rdb, err := tredis.GetRedis(8)
	//建立redis连接
	defer rdb.Close()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Error("redis连接错误")
		return err
	}
	//写入消息失败的回调函数
	addBack := func(k string, v string) {
		logmw.LogWithRequest("Favorite", F.C).Error("kafka写入失败")
		rdb.ZRem(Rctx, userIdStr, videoId)
		exi, _ := rdb.Exists(Rctx, videoIdStr).Result()
		if exi > 0 {
			rdb.Decr(Rctx, videoIdStr)
		}
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":Kafka写入失败"))
	}
	delBack := func(k string, v string) {
		logmw.LogWithRequest("Favorite", F.C).Warn("kafka写入失败")
		rdb.ZAdd(Rctx, userIdStr, redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: v,
		})
		exi, _ := rdb.Exists(Rctx, videoIdStr).Result()
		if exi > 0 {
			rdb.Incr(Rctx, videoIdStr)
		}
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":Kafka写入失败"))
	}
	//判断缓存中key是否存在
	exists, err := rdb.Exists(Rctx, userIdStr).Result()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Error("redis连接错误")
		return err
	}
	if exists < 1 {
		//缓存中不存在,更新缓存
		err = LoadFavoriteToRides(userId, rdb, Rctx)
		if err != nil {
			F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
			logmw.LogWithRequestErr("Favorite", F.C, err).Error("redis更新错误")
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
		Now := time.Now().Unix()
		n, err = rdb.ZAdd(Rctx, userIdStr, redis.Z{
			Score:  float64(Now),
			Member: videoIdStr,
		}).Result()
		if err != nil {
			F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
			logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis写入失败")
			return err
		}
		if n == 0 {
			logmw.LogWithRequest("Favorite", F.C).Debug("重复点赞")
			return nil
		}
		//更新视频被点赞数量
		exists, _ = rdb.Exists(Rctx, videoIdStr).Result()
		if exists > 0 {
			rdb.Incr(Rctx, videoIdStr)
		}

		//更新redis成功，向消息队列发送key:UserId  value:videoId+time.now().unix(),为方便直接使用字符串拼接，追求性能可以使用结构体将结构体序列化
		value := videoIdStr + " " + strconv.FormatInt(Now, 10)
		kafka.FavoriteMq.WriteMsg(userIdStr, value, addBack)
		return nil
	} else {
		//取消点赞
		n, err = rdb.ZRem(Rctx, userIdStr, videoId).Result()
		if err != nil {
			F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
			logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis删除失败")
			return err
		}
		if n == 0 {
			logmw.LogWithRequest("Favorite", F.C).Debug("重复取消赞")
			return nil
		}
		//更新视频被点赞数量
		exists, _ = rdb.Exists(Rctx, videoIdStr).Result()
		if exists > 0 {
			rdb.Decr(Rctx, videoIdStr)
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
	var zs []redis.Z
	for _, val := range VideoIds {
		zs = append(zs, redis.Z{Score: float64(val.CreatedAt), Member: val.VideoId})
	}
	UserIdStr := strconv.FormatInt(UserId, 10)
	//添加占位，防止点赞全部取消后key被删除，此时数据库还未更新
	err = rdb.ZAdd(Rctx, UserIdStr, redis.Z{
		Score:  -1,
		Member: -1,
	}).Err()
	err = rdb.ZAdd(Rctx, UserIdStr, zs...).Err()
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

// GetFavouriteList 获取点赞视频列表
func (F *FavoriteImpl) GetFavouriteList(userId int64, curId int64) ([]*model.Video, error) {
	userIdStr := strconv.FormatInt(userId, 10)
	VideoServer := &VideoServiceImpl{
		C: F.C,
	}
	var Rctx = context.Background()
	rdb, err := tredis.GetRedis(8)
	//建立redis连接
	defer rdb.Close()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis连接错误")
		return nil, err
	}
	//判断缓存中key是否存在
	exists, err := rdb.Exists(Rctx, userIdStr).Result()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis连接错误")
		return nil, err
	}
	if exists < 1 {
		//缓存中不存在,更新缓存
		err = LoadFavoriteToRides(userId, rdb, Rctx)
		if err != nil {
			F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
			logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis更新错误")
			return nil, err
		}
	} else {
		//缓存中存在，更新缓存时间
		if ttl, _ := rdb.TTL(Rctx, userIdStr).Result(); ttl < constant.Favorite_UserId_DefaultTime/3 {
			rdb.Expire(Rctx, userIdStr, constant.Favorite_UserId_DefaultTime/3)
		}
	}
	//获取所有点赞列表
	var VideoIds []int64
	vals, _ := rdb.ZRevRange(Rctx, userIdStr, 0, -1).Result()
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
				logmw.LogWithRequestErr("favorite", F.C, err).WithField("videoId:", id).Error("视频获取出错")
				wg.Done()
				return
			}
			VideoList[a] = Video
			wg.Done()
		}(i, VideoId)
	}
	wg.Wait()
	return VideoList, nil
}

// IsFavorite 判断是否点赞该视频
func (F *FavoriteImpl) IsFavorite(videoId int64, userId int64) (bool, error) {
	userIdStr := strconv.FormatInt(userId, 10)
	vidoeIdStr := strconv.FormatInt(videoId, 10)
	var Rctx = context.Background()
	rdb, err := tredis.GetRedis(8)
	exists, err := rdb.Exists(Rctx, userIdStr).Result()
	if err != nil {
		logmw.LogWithRequestErr("Favorite", F.C, err).Warn("redis连接错误")
	}
	//缓存中不存在去数据库中查
	if exists < 1 {
		return dao.ExistsFav(userId, videoId)
	}
	val, _ := rdb.ZScore(Rctx, userIdStr, vidoeIdStr).Result()
	if val == float64(0) {
		return false, nil
	}
	return true, nil
}

// FavouriteVideoCount 根据userId获取这个用户点赞视频数量
func (F *FavoriteImpl) FavouriteVideoCount(userId int64) (n int64, err error) {
	userIdStr := strconv.FormatInt(userId, 10)
	var Rctx = context.Background()
	rdb, err := tredis.GetRedis(8)
	//建立redis连接
	defer rdb.Close()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis连接错误")
		return 0, err
	}
	//判断缓存中key是否存在
	exists, err := rdb.Exists(Rctx, userIdStr).Result()
	if err != nil {
		F.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("Favorite", F.C, err).Debug("redis连接错误")
		return 0, err
	}
	if exists > 0 {
		n, err = rdb.ZCard(Rctx, userIdStr).Result()
		n = n - 1 //减去占位符
		return n, err
	}
	//缓存中不存在，查数据库
	n, err = dao.GetUserFavorCount(userId)
	return
}

// FavouriteCount 根据当前视频id获取当前视频点赞数量。
func (F *FavoriteImpl) FavouriteCount(videoId int64) (int64, error) {
	//缓存中 videoId->count
	videoIdStr := strconv.FormatInt(videoId, 10)
	var Rctx = context.Background()
	rdb, err := tredis.GetRedis(8)
	//建立redis连接
	defer rdb.Close()
	if err != nil {
		logmw.LogWithRequestErr("Favorite", F.C, err).Error("redis连接错误")
		return 0, err
	}
	//判断缓存中key是否存在
	exists, err := rdb.Exists(Rctx, videoIdStr).Result()
	if err != nil {
		logmw.LogWithRequestErr("Favorite", F.C, err).Error("redis连接错误")
	}
	if exists < 1 {
		//缓存中不存在,更新缓存
		count, err := dao.GetVideoFavorCount(videoId)
		if err != nil {
			logmw.LogWithRequestErr("Favorite", F.C, err).Warn("mysql获取视频点赞数量出错")
			return count, err
		}
		rdb.Set(Rctx, videoIdStr, count, constant.Favorite_UserId_DefaultTime)
		return count, nil
	} else {
		//缓存中存在，更新缓存时间
		if ttl, _ := rdb.TTL(Rctx, videoIdStr).Result(); ttl < constant.Favorite_UserId_DefaultTime/3 {
			rdb.Expire(Rctx, videoIdStr, constant.Favorite_UserId_DefaultTime/3)
		}
	}
	n, err := rdb.Get(Rctx, videoIdStr).Int64()
	if err != nil {
		logmw.LogWithRequestErr("Favorite", F.C, err).Warn("redis获取视频点赞数量出错")
		n, err = dao.GetVideoFavorCount(videoId)
		if err != nil {
			logmw.LogWithRequestErr("Favorite", F.C, err).Warn("mysql获取视频点赞数量出错")
			return n, err
		}
		return n, nil
	}
	return n, nil
}
