package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/biz/model"
	"tiktok/biz/service"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	"time"
)

type FeedResp struct {
	model.BaseResp
	VideoList []*model.Video `json:"video_list"`
	NextTime  int64          `json:"next_time"`
}

type VideoListResp struct {
	*model.BaseResp
	VideoList []*model.Video `json:"video_list"`
}

// Feed 视频Feed流：
//
// @Router /douyin/feed [get]
// @Summary 视频Feed流
// @Schemes
// @Description 支持所有用户刷抖音，视频按投稿时间倒序推出
// @Tags Video
// @Accept json
// @Produce json
// @Param latest_time query string false "上次请求的时间戳"
// @Param user_id query int false "用户ID"
// @Success		200
// @Success 	200 {object} model.BaseResp
// @Failure 	400 {object} model.BaseResp
// @Failure     500  {object}  model.BaseResp
func Feed(c *gin.Context) {
	// 可选参数：latestTime
	latestTime := Timestamp2Time(c.Query("latest_time"))
	// 可选参数 userId
	userId, err := getUserId(c)
	if err != nil {
		userId = int64(0)
	}

	videoList, newTime, err := service.NewVideoServiceImpl(c).Feed(latestTime, userId)
	if err != nil {
		utils.LogWithRequestId(c, "video", err).Debug("Service异常")
		c.JSON(http.StatusInternalServerError, errno.VideoIsNotExistErr)
		return
	}

	// 业务有错误，不响应
	if c.IsAborted() {
		return
	}

	c.JSON(
		http.StatusOK,
		FeedResp{
			BaseResp: model.BaseResp{
				Code: 0,
				Msg:  "success",
			},
			VideoList: videoList,
			NextTime:  newTime.Unix(),
		})
}

// UpVideo 视频投稿：支持登录用户自己拍视频投稿
// @Router /douyin/publish/action/ [post]
// @Summary 视频投稿
// @Schemes
// @Description 支持登录用户自己拍视频投稿
// @Tags Video
// @Accept json
// @Produce json
func UpVideo(c *gin.Context) {
	// 获取 userId
	userId := c.GetInt64(constant.UserId)
	//userId, err := getUserId(c)

	// 提取POST参数
	req := &model.VideoUploadReq{}

	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.BuildBindResp(err))
		utils.LogParamError(c, err)
		return
	}

	req.Data, req.File, err = c.Request.FormFile("data")
	if err != nil {
		utils.LogParamError(c, err)
		return
	}

	// 如果视频过大100Mb
	if req.File.Size > (100 << 20) {
		c.JSON(http.StatusOK, errno.VideoUploadFailedErr)
		return
	}

	// 上传视频
	err = service.NewVideoServiceImpl(c).Publish(req.Data, req.File, userId, req.Title)
	if err != nil {
		c.JSON(http.StatusOK, errno.VideoUploadFailedErr)
		return
	}

	c.JSON(http.StatusOK, model.BaseResp{
		Code: 0,
		Msg:  "success",
	})
}

// VideoList 获取用户发布的视频列表
// @Router /douyin/video/list/ [get]
// @Summary 获取用户发布的视频列表
// @Schemes
// @Description 获取用户发布的视频列表
// @Tags Video
// @Accept json
// @Produce json
// @Param user_id query int64 true "用户ID"
// @Success		200
// @Failure 	401 {body}  errno.Errno "未登录"
func VideoList(c *gin.Context) {
	// 获取本用户ID
	userId := c.GetInt64(constant.UserId)

	// 获取目标用户ID
	targetUserId, err := getUserId(c)
	if err != nil {
		utils.LogParamError(c, err)
		return
	}

	// 获取视频列表
	videoList, err := service.NewVideoServiceImpl(c).GetVideoListById(targetUserId, userId)
	if err != nil {
		utils.LogWithRequestId(c, "video", err).Debug("Service异常")
		c.JSON(http.StatusInternalServerError, errno.VideoServiceErr)
		return
	}

	// 业务有错误，不响应
	if c.IsAborted() {
		return
	}

	c.JSON(http.StatusOK, VideoListResp{
		BaseResp:  model.BuildBaseResp(err),
		VideoList: videoList,
	})
}

// Timestamp2Time 时间戳转换为time.Time
func Timestamp2Time(str string) time.Time {
	if len(str) == 0 {
		return time.Now()
	}
	// 解析时间
	timestamp, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		// 处理转换错误
		return time.Now()
	}
	// 使用时间戳创建 time.Time 对象
	return time.Unix(timestamp, 0)
}
