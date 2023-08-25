package service

import (
	"context"
	"encoding/json"
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
	"tiktok/pkg/utils"
	"time"
)

type CommentService interface {
	// GetCommCountFromVId 根据视频id获取评论数量
	GetCommCountFromVId(id int64) (int64, error)
	AddComment(userId int64, videoId int64, text string) (*model.Comment, error)
	DelComment(VideoId, CommId int64) error
	GetCommentList(videoId int64) ([]*model.Comment, error)
}
type CommentServiceImpl struct {
	C *gin.Context
	UserService
}

// 评论使用旁路缓存策略，保证高一致性
// 添加评论
func (comm *CommentServiceImpl) AddComment(userId int64, videoId int64, text string) (*model.Comment, error) {
	//判断视频是否存在
	VideoS := &VideoServiceImpl{C: comm.C}
	UserS := &UserServiceImpl{C: comm.C}
	if _, err := VideoS.GetVideoById(videoId, userId); err != nil {
		comm.C.AbortWithStatusJSON(http.StatusBadRequest, errno.NewErrno(errno.VideoIsNotExistErrCode, "视频不存在"))
		logmw.LogWithRequestErr("comment", comm.C, err).Warn("视频不存在")
		return nil, err
	}
	back := func(key string, value string) {
		comm.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":kafka写入失败"))
		logmw.LogWithRequest("comment", comm.C).Error("kafka写入失败")
	}
	createAt := time.Now().Unix()
	commId := utils.UUidToInt64ID()
	user := UserS.GetUserByUserId(userId)
	commModel := &model.CommToJson{
		CommId:   commId,
		CreateAt: createAt,
		UserId:   userId,
		VideoId:  videoId,
		Content:  text,
	}
	commjson, _ := json.Marshal(commModel)
	//将评论信息写入kafka即可返回评论成功
	commReq := &model.Comment{
		Id:        commId,
		CreatedAt: createAt,
		Content:   text,
		User:      user,
	}
	kafka.CommonMq.WriteMsg("Add", string(commjson), back)
	return commReq, nil
}

// 删除评论
func (comm *CommentServiceImpl) DelComment(VideoId, CommId int64) error {
	VideoIdStr := strconv.FormatInt(VideoId, 10)
	CommIdStr := strconv.FormatInt(CommId, 10)
	//删除缓存中数据，然后发送kafka消息
	Rctx := context.Background()
	rdb, err := tredis.GetRedis(9)
	if err != nil {
		comm.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("comment", comm.C, err).Error("redis连接错误")
		return err
	}
	exi1, _ := rdb.Exists(Rctx, VideoIdStr).Result()
	if exi1 > 0 {
		rdb.ZRem(Rctx, VideoIdStr, CommId)
	}
	exi2, _ := rdb.Exists(Rctx, CommIdStr).Result()
	if exi2 > 0 {
		rdb.Del(Rctx, CommIdStr)
	}
	back := func(string1 string, string2 string) {
		rdb.Del(Rctx, VideoIdStr)
		comm.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":kafka写入失败"))
		logmw.LogWithRequest("comment", comm.C).Error("kafka写入失败")
	}
	kafka.CommonMq.WriteMsg("Del", CommIdStr, back)
	return nil
}

// 获取评论列表
func (comm *CommentServiceImpl) GetCommentList(videoId int64) ([]*model.Comment, error) {
	videoIdStr := strconv.FormatInt(videoId, 10)
	userS := UserServiceImpl{C: comm.C}
	var Rctx = context.Background()
	rdb, err := tredis.GetRedis(8)
	//建立redis连接
	defer rdb.Close()
	if err != nil {
		comm.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("comment", comm.C, err).Debug("redis连接错误")
		return nil, err
	}
	exists, err := rdb.Exists(Rctx, videoIdStr).Result()
	if err != nil {
		comm.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr.AppendMsg(":RedisErr"))
		logmw.LogWithRequestErr("comment", comm.C, err).Debug("redis连接错误")
		return nil, err
	}
	wg := sync.WaitGroup{}
	if exists < 1 {
		//缓存中不存在,查询数据库更新缓存
		comms, err := dao.GetCommList(videoId)
		if err != nil {
			comm.C.AbortWithStatusJSON(http.StatusInternalServerError, errno.ServiceErr)
			logmw.LogWithRequestErr("comment", comm.C, err).Debug("数据库连接错误")
			return nil, err
		}
		go func(wg *sync.WaitGroup) {
			wg.Add(1)
			defer wg.Done()
			err := LoadCommIdsToRedis(videoId, comms, rdb, Rctx)
			if err != nil {
				logmw.LogWithRequestErr("comment", comm.C, err).Warn("缓存载入失败")
			}
		}(&wg)
		for i, val := range comms {
			wg.Add(1)
			go func(i int, v *model.Comment) {
				comms[i].User = userS.GetUserByUserId(comms[i].UserId)
				wg.Done()
			}(i, val)
		}
		wg.Wait()
		return comms, nil
	} else {
		//缓存中存在
		//按照时间戳逆序获取评论id
		commIds, err := rdb.ZRevRangeWithScores(Rctx, videoIdStr, 0, -1).Result()
		if err != nil {
			//降级
			logmw.LogWithRequestErr("comment", comm.C, err).Warn("从缓存中获取评论失败")
			return dao.GetCommList(videoId)
		}
		//组装评论
		commList := make([]*model.Comment, len(commIds))
		for i, val := range commIds {
			go func(n int, v redis.Z) {
				wg.Add(1)
				defer wg.Done()
				commList[n].Id, _ = strconv.ParseInt(v.Member.(string), 10, 64)
				commList[n].CreatedAt = int64(v.Score)
				var Info CommInfo
				err := rdb.HGetAll(Rctx, v.Member.(string)).Scan(&Info)
				if err != nil {
					logmw.LogWithRequestErr("comment", comm.C, err).Warn("从缓存中获取评论失败")
					//获取失败，降级从数据库中获取
					commList[n], err = dao.Comm(commList[n].Id)
					if err != nil {
						logmw.LogWithRequestErr("comment", comm.C, err).Warn("缓存获取失败后从数据库获取评论失败")
					}
				}
			}(i, val)
		}
		wg.Wait()
		return commList, nil
	}
}

// 根据视频id获取评论数量
func (comm *CommentServiceImpl) GetCommCountFromVId(id int64) (int64, error) {
	count, err := dao.GetCommCount(id)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// 从数据库中将指定评论内容加载到redis中
func LoadCommToRedis(rdb redis.Client, Rctx context.Context, commId int64) error {
	comm, err := dao.Comm(commId)
	if err != nil {
		return err
	}
	err = rdb.HMSet(Rctx, strconv.FormatInt(comm.Id, 10), CommInfo{
		UserId:  comm.UserId,
		Content: comm.Content,
	}, constant.Comment_CommId_DefaultTime).Err()
	return err
}

type CommInfo struct {
	UserId  int64  `redis:"userid"`
	Content string `redis:"content"`
}

// 将视频的评论加载到redis中
func LoadCommIdsToRedis(VideoId int64, CommIds []*model.Comment, rdb *redis.Client, Rctx context.Context) error {
	VideoIdStr := strconv.FormatInt(VideoId, 10)
	var zs []redis.Z
	for _, val := range CommIds {
		zs = append(zs, redis.Z{Score: float64(val.CreatedAt), Member: val.Id})
		rdb.HMSet(Rctx, strconv.FormatInt(val.Id, 10), CommInfo{
			UserId:  val.UserId,
			Content: val.Content,
		}, constant.Comment_CommId_DefaultTime)
	}
	zs = append(zs, redis.Z{
		Score:  -1,
		Member: -1,
	})
	err := rdb.ZAdd(Rctx, VideoIdStr, zs...).Err()
	if err != nil {
		return err
	}
	//设置过期时间，默认3天
	err = rdb.Expire(Rctx, VideoIdStr, constant.Favorite_UserId_DefaultTime).Err()
	if err != nil {
		return err
	}
	return nil
}
