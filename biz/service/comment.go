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
	AddComment(userId int64, videoId int64, text string) *model.Comment
	DelComment(userId, CommId int64)
	GetCommentList(videoId int64) ([]*model.Comment, error)
}
type CommentServiceImpl struct {
	C *gin.Context
	UserService
}

// 添加评论
func (comm *CommentServiceImpl) AddComment(userId int64, videoId int64, text string) *model.Comment {
	//判断视频是否存在
	VideoS := &VideoServiceImpl{C: comm.C}
	UserS := &UserServiceImpl{ctx: comm.C}
	if _, err := VideoS.GetVideoById(videoId, userId); err != nil {
		utils.LogWithRequestId(comm.C, "Comment", err)
		utils.LogBizErr(comm.C, errno.VideoIsNotExistErr, http.StatusOK, "视频不存在")
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
	kafka.CommonMq.WriteMsg("Add", string(commjson), func(key string, value string) {
		utils.LogBizErr(comm.C, errno.CommentActionErr, http.StatusOK, "kafka写入失败")
	})
	return commReq
}

// 删除评论
func (comm *CommentServiceImpl) DelComment(UserId, CommId int64) {
	var Rctx = context.Background()
	rdb, err := tredis.CommR.GetCommRedis()
	if err != nil {
		utils.LogDB(comm.C, errno.Service)
	}
	//判断是否有权限删除
	err = tredis.LoadIfNotExists(CommId, rdb, tredis.LoadCommToRedis)
	if err != nil {
		utils.LogBizErr(comm.C, errno.Update, http.StatusOK, "获取评论失败")
	}
	Info := model.CommInfo{}
	_ = rdb.HGetAll(Rctx, strconv.FormatInt(CommId, 10)).Scan(&Info)
	if Info.UserId != UserId {
		utils.LogBizErr(comm.C, errno.CommentActionErr, http.StatusOK, "没有删除权限")
	}
	kafka.CommonMq.WriteMsg("Del", strconv.FormatInt(CommId, 10), func(key string, value string) {
		utils.LogBizErr(comm.C, errno.CommentActionErr, http.StatusOK, "kafka写入失败")
	})
}

// 获取评论列表
func (comm *CommentServiceImpl) GetCommentList(videoId int64) []*model.Comment {
	var wg sync.WaitGroup
	//判断视频是否存在
	UserS := &UserServiceImpl{ctx: comm.C}
	VideoS := &VideoServiceImpl{C: comm.C}
	if _, err := VideoS.GetVideoById(videoId, comm.C.GetInt64(constant.UserId)); err != nil {
		utils.LogBizErr(comm.C, errno.VideoIsNotExistErr, http.StatusOK, "视频不存在")
		return nil
	}
	var Rctx = context.Background()
	rdb, err := tredis.CommR.GetCommRedis()
	if err != nil {
		utils.LogDB(comm.C, errno.Service)
		return nil
	}
	err = tredis.LoadIfNotExists(videoId, rdb, tredis.LoadCommsToRedis)
	if err != nil {
		utils.LogBizErr(comm.C, errno.Update, http.StatusOK, "将评论列表加载失败")
		return nil
	}
	//按照时间戳逆序获取评论id
	var commList []*model.Comment
	commIds, err := rdb.ZRevRangeWithScores(Rctx, strconv.FormatInt(videoId, 10), 0, -1).Result()
	if err != nil {
		//降级,从mysql中查询
		utils.LogWithRequestId(comm.C, "comment", err).Warn("从缓存中获取评论失败")
		commList, err = dao.GetCommList(videoId)
		for _, val := range commList {
			go func(comment *model.Comment) {
				wg.Add(1)
				defer wg.Done()
				comment.User = UserS.GetUserByUserId(comment.UserId)
			}(val)
		}
		wg.Wait()
		return commList
	}
	//组装评论
	commList = make([]*model.Comment, len(commIds))
	for i, val := range commIds {
		go func(n int, v redis.Z) {
			wg.Add(1)
			defer wg.Done()
			commList[n].Id, _ = strconv.ParseInt(v.Member.(string), 10, 64)
			commList[n].CreatedAt = int64(v.Score)
			var Info model.CommInfo
			err := rdb.HGetAll(Rctx, v.Member.(string)).Scan(&Info)
			if err != nil {
				utils.LogWithRequestId(comm.C, "comment", err).Debug("从缓存中获取失败")
				//获取失败，降级从数据库中获取
				commList[n], err = dao.Comm(commList[n].Id)
				if err != nil {
					utils.LogWithRequestId(comm.C, "comment", err).Error("从mysql中获取评论失败")
				}
				commList[n].User = UserS.GetUserByUserId(commList[n].UserId)
			}
			commList[n].Content = Info.Content
			commList[n].User = UserS.GetUserByUserId(Info.UserId)
		}(i, val)
	}
	wg.Wait()
	return commList
}

// 根据视频id获取评论数量
func (comm *CommentServiceImpl) GetCommCountFromVId(id int64) (int64, error) {
	count, err := dao.GetCommCount(id)
	if err != nil {
		return 0, err
	}
	return count, nil
}
