package service

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strconv"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/oss"
	"tiktok/biz/model"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
	"time"
)

type VideoService interface {
	//Feed 传入时间戳,当前用户Id，返回视频切片 和 返回的视频切片中的最早时间
	Feed(lastTime time.Time, userId int64) ([]*model.Video, time.Time, error)
	//GetVideoById 根据视频id和用户id获取video
	GetVideoById(videoId int64, userId int64) (*model.Video, error)
	//Publish 上传视频
	Publish(file multipart.File, fileHeader *multipart.FileHeader, userId int64, title string) error
	//GetVideoListById 当前用户 (userId) 获取目标用户(targetId)发布的视频
	GetVideoListById(targetId int64, userId int64) ([]*model.Video, error)
}

func NewVideoServiceImpl(c *gin.Context) *VideoServiceImpl {
	return &VideoServiceImpl{
		ctx: c,
	}
}

// Feed 根据时间戳和用户ID，获取用户的动态消息流视频列表及最早时间。
// 该函数会根据传入的时间戳和用户ID，检索在指定时间之后用户的动态消息流视频列表，并返回这些视频对象的切片及其中最早发布的时间。
//
// 参数:
//   lastTime: 上一次获取时间的时间戳
//   userId: 用户的ID
//
// 返回值:
//   videos: 用户的动态消息流视频列表，类型为 []*model.Video
//   earliestTime: 返回视频切片中最早发布的时间
//   error: 如果在处理过程中出现错误，会返回相应的错误对象
func (v VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]*model.Video, time.Time, error) {
	limitNum := 30

	// 获取视频列表
	videoList, err := dao.NewVideo().GetVideoListByLastTime(lastTime, limitNum)
	if err != nil {
		// 数据库查询出错
		utils.LogDB(v.ctx, err)
		return nil, lastTime, err
	}
	if len(videoList) <= 0 {
		// 空视频列表，由Handler端处理
		return videoList, lastTime, nil
	}

	utils.LogWithData("videoList", videoList).Debug("Feed方法：获取视频列表")

	// 完善Author和Favor信息
	err = v.completeVideoList(videoList, userId)
	if err != nil {
		utils.LogDB(v.ctx, err)
		return nil, lastTime, err
	}

	// 新的lastTime为列表中最新视频的Created_at表示的时间
	newLastTime := time.Unix(videoList[len(videoList)-1].CreatedAt, 0)
	return videoList, newLastTime, nil
}

// GetVideoById 根据视频ID和用户ID从数据库中获取视频信息。
//
// 参数:
// - videoId: 要查询的视频ID。
// - userId: 查询用户的ID，用于权限验证。
//
// 返回值:
// - *model.Video: 表示视频的模型对象，包含视频的详细信息。
// - error: 如果查询过程中出现错误，将返回一个非空的错误对象，否则为nil。
//
// 注意事项:
// - 如果未找到与给定视频ID相对应的视频，返回值中的模型对象将为nil。
// - 如果用户无权限访问该视频，同样会返回nil的模型对象。
// - 若出现任何数据库查询错误，将返回相应的错误信息。
func (v VideoServiceImpl) GetVideoById(videoId int64, userId int64) (*model.Video, error) {
	// VideoDAO 获取原始Video对象
	tmpVideo := &model.Video{Id: videoId}
	targetVideo, err := dao.NewVideo().GetVideoById(tmpVideo)
	if err != nil {
		// 数据库查询出错
		utils.LogDB(v.ctx, err)
		return nil, err
	}
	// 查询Author
	videoAuthor, err := dao.MustGetUserById(targetVideo.AuthorId)
	if err != nil {
		utils.LogDB(v.ctx, err)
		return nil, err
	}
	targetVideo.Author = videoAuthor

	// 查询Favor
	isFavor, err := dao.ExistsFav(userId, videoId)
	if err != nil {
		utils.LogDB(v.ctx, err)
		return nil, err
	}
	targetVideo.IsFavorite = isFavor

	return targetVideo, nil
}

// Publish 将上传的视频文件发布为用户的动态消息。
// 该函数会将视频文件上传到云存储，获取视频链接和封面链接，并创建一个新的视频对象。
// 参数:
//   file: 上传的视频文件
//   fileHeader: 上传文件的文件头信息
//   userId: 用户的ID
//   title: 视频标题
// 返回值:
//   error: 如果在处理过程中出现错误，会返回相应的错误对象
func (v VideoServiceImpl) Publish(file multipart.File, fileHeader *multipart.FileHeader, userId int64, title string) error {

	// 上传视频文件到腾讯云cos, 获取视频链接videoUrl、封面链接coverUrl
	videoUrl, coverUrl, err := oss.UpLoad(file, strconv.FormatInt(userId, 10), fileHeader)
	if err != nil {
		// 调用cos服务出错
		utils.LogWithRequestId(v.ctx, "video", err).Debug("Publish方法：调用cos服务出错")
		return err
	}
	newVideo := model.Video{
		//Id:            0, // 不提供Id
		AuthorId:      userId,
		Title:         title,
		PlayUrl:       videoUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		// 创建User对象
		Author: &model.User{
			Id: userId,
		},
	}

	err = dao.NewVideo().Create(&newVideo)
	if err != nil {
		// 数据库新增出错
		utils.LogDB(v.ctx, err)
		return err
	}
	return nil
}

func (v VideoServiceImpl) completeVideoList(videoList []*model.Video, userId int64) error {
	VideoIdList := make([]int64, len(videoList))
	UserIdList := make([]int64, 0)
	userIdMap := make(map[int64]*model.User)
	uniqueIdMap := make(map[int64]bool)

	for idx, video := range videoList {
		VideoIdList[idx] = video.Id
		if !uniqueIdMap[video.AuthorId] {
			uniqueIdMap[video.AuthorId] = true
			UserIdList = append(UserIdList, video.AuthorId)
		}
	}

	//--------------------- Author ---------------------

	// 查询Author
	authorList, err := dao.MustGetUsersByIds(UserIdList)
	if err != nil {
		utils.LogDB(v.ctx, err)
		return err
	}

	// 根据User列表构造Map
	for _, user := range authorList {
		userIdMap[user.Id] = user
	}

	// 将User列表赋值给Video列表
	for _, video := range videoList {
		video.Author = userIdMap[video.AuthorId] // 从Map中获取User
	}

	//------------------ favor ---------------------

	// 查询Favor
	favorList, err := dao.CheckFavorListByUserAndVideo(VideoIdList, userId)
	if err != nil {
		utils.LogDB(v.ctx, err)
		return err
	}
	// 如果favor调用出现错误，就都为false
	if len(favorList) != len(videoList) {
		return errno.FavoriteActionErr
	}

	for i, video := range videoList {
		video.IsFavorite = favorList[i]
	}

	return nil
}

// GetVideoListById 当前用户 (userId) 获取目标用户(targetId)发布的视频
// 参数：
// - targetId: 目标用户的ID，表示要获取视频列表的用户。
// - userId: 当前用户的ID，表示正在访问视频列表的用户。
// 返回值：
// - 一个指向 model.Video 切片的指针，表示目标用户发布的视频列表。
// - 如果在查询数据库时出现问题，则返回错误。
func (v VideoServiceImpl) GetVideoListById(targetId int64, userId int64) ([]*model.Video, error) {
	targetVideoList, err := dao.NewVideo().GetVideoListByAuthor(targetId)
	if err != nil {
		utils.LogDB(v.ctx, err)
		return nil, err
	}

	err = v.completeVideoList(targetVideoList, userId)
	if err != nil {
		return nil, err
	}

	return targetVideoList, nil
}

type VideoServiceImpl struct {
	ctx *gin.Context
}
