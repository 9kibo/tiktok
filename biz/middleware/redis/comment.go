package redis

import (
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"tiktok/biz/dao"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
)

// 复用全局连接
type CommRedis struct {
	Comment *redis.Client
	sync.Once
}

var CommMutex = utils.NewCacheGuard()

var CommR = CommRedis{}

func (Comm *CommRedis) GetCommRedis() (*redis.Client, error) {
	var err error
	Comm.Do(func() {
		client, err := GetRedis(9)
		if err != nil {
			utils.Log("redis").WithField("err:", err).Error("redis连接失败")
		}
		Comm.Comment = client
	})
	err = Comm.Comment.Ping(Ctx).Err()
	if err != nil {
		utils.Log("redis").WithField("err:", err).Error("redis连接失败")
		return nil, err
	}
	return Comm.Comment, nil
}

// videoId ->ZSet(score(createAt),Member:CommId)
// CommId  ->Hash model.CommInfo{
//		userId
//		content
//	}

// LoadCommsToRedis 将评论列表加载到redis
func LoadCommsToRedis(key int64, rdb *redis.Client) error {
	CommIds, err := dao.GetCommList(key)
	VideoIdStr := strconv.FormatInt(key, 10)
	var zs []redis.Z
	for _, val := range CommIds {
		zs = append(zs, redis.Z{Score: float64(val.CreatedAt), Member: val.Id})
		rdb.HMSet(Ctx, strconv.FormatInt(val.Id, 10), model.CommInfo{
			UserId:  val.UserId,
			Content: val.Content,
		}, constant.Comment_CommId_DefaultTime)
	}
	zs = append(zs, redis.Z{
		Score:  -1,
		Member: -1,
	})
	err = rdb.ZAdd(Ctx, VideoIdStr, zs...).Err()
	if err != nil {
		return err
	}
	//设置过期时间
	rdb.Expire(Ctx, VideoIdStr, constant.Favorite_UserId_DefaultTime)
	return nil
}

// LoadCommToRedis 将单个评论加载到redis
func LoadCommToRedis(key int64, rdb *redis.Client) error {
	comm, err := dao.Comm(key)
	if err != nil {
		return err
	}
	err = rdb.HMSet(Ctx, strconv.FormatInt(comm.Id, 10), model.CommInfo{
		UserId:  comm.UserId,
		Content: comm.Content,
	}, constant.Comment_CommId_DefaultTime).Err()
	return err
}
