package redis

import (
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"tiktok/biz/dao"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
)

// 复用全局连接
type FavRedis struct {
	Fav *redis.Client
	sync.Once
}

var FavMutex = utils.NewCacheGuard()

var FavR = FavRedis{}

func (F *FavRedis) GetFavRedis() (*redis.Client, error) {
	var err error
	F.Do(func() {
		client, err := GetRedis(8)
		if err != nil {
			utils.Log("redis").WithField("err:", err).Error("redis连接失败")
		}
		F.Fav = client
	})
	err = F.Fav.Ping(Ctx).Err()
	if err != nil {
		utils.Log("redis").WithField("err:", err).Error("redis连接失败")
		return nil, err
	}
	return F.Fav, nil
}

func UpdateFavRedis(rdb *redis.Client, UserId int64, VideoId int64, ActionType int32, Creat int64, Rid string) {
	//1点赞 2取消点赞
	err := LoadIfNotExists(UserId, rdb, LoadFavoriteToRides)
	if err != nil {
		utils.LogWithRidString("FavRedis", Rid, err).Error("更新缓存失败")
	}
	//更新缓存
	if ActionType == 1 {
		err := rdb.ZAdd(Ctx, strconv.FormatInt(UserId, 10), redis.Z{
			Score:  float64(Creat),
			Member: VideoId,
		}).Err()
		if err != nil {
			utils.LogWithRidString("FavRedis", Rid, err).Error("更新缓存失败")
		}
		//更新视频被点赞数量
		exists, _ := rdb.Exists(Ctx, strconv.FormatInt(VideoId, 10)).Result()
		if exists > 0 {
			rdb.Incr(Ctx, strconv.FormatInt(VideoId, 10))
		}
	} else {
		err = rdb.ZRem(Ctx, strconv.FormatInt(UserId, 10), strconv.FormatInt(VideoId, 10)).Err()
		if err != nil {
			utils.LogWithRidString("FavRedis", Rid, err).Error("更新缓存失败")
		}
		//更新视频被点赞数量
		exists, _ := rdb.Exists(Ctx, strconv.FormatInt(VideoId, 10)).Result()
		if exists > 0 {
			rdb.Decr(Ctx, strconv.FormatInt(VideoId, 10))
		}
	}

}
func LoadFavoriteToRides(UserId int64, rdb *redis.Client) error {
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
	zs = append(zs, redis.Z{
		Score:  -1,
		Member: -1,
	})
	err = rdb.ZAdd(Ctx, UserIdStr, zs...).Err()
	if err != nil {
		return err
	}
	//设置过期时间
	err = rdb.Expire(Ctx, UserIdStr, constant.Favorite_UserId_DefaultTime).Err()
	if err != nil {
		return err
	}
	return nil
}

// 将视频被点赞数加载到redis中
func LoadVideoFavCount(VideoId int64, rdb *redis.Client) error {
	count, err := dao.GetVideoFavorCount(VideoId)
	if err != nil {
		return err
	}
	err = rdb.Set(Ctx, strconv.FormatInt(VideoId, 10), count, constant.Favorite_UserId_DefaultTime).Err()
	return err
}
