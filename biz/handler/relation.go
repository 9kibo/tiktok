// Package handler
//  1. service方法写就不返回, 读也仅仅是返回结果
//  2. 在service中调用ctx.AbortWithStatusJSON()响应错误, 在handler中调用ctx.IsAborted()判断是否出错，出错就直接返回
//     (原因: 在service控制响应)
package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/biz/model"
	"tiktok/biz/service"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

// RelationAction
//
// @Router /douyin/relation/action [post]
// @Summary 关注操作
// @Schemes http
// @Description 关注或者取消关注用户
// @Tags Relation
// @Accept       json
// @Produce      json
// @Param user_id query int64 true "关注者"
// @Param to_user_id query int64 true "被关注者"
// @Success      200
// @Failure      400  {body}  errno.Errno "参数不正确"
// @Failure      403  {body}  errno.Errno "未登录"
// @Failure      500  {body}  errno.Errno "系统错误"
func RelationAction(ctx *gin.Context) {
	//bind arg and validate arg
	req := &model.FollowActionReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		utils.LogParamError(ctx, err)
		return
	}

	//do service
	service.NewFollowService(ctx).FollowAction(req)

	//service has error (already send error resp)
	if ctx.IsAborted() {
		return
	}

	//service hasn't error, send success resp
	respOfUpdate(ctx)
}

// FollowListResp 关注/粉丝/好友列表
type FollowListResp struct {
	errno.Errno
	//用户列表
	UserList []*model.User `json:"user_list"`
}

// GetFollowingList
//
// @Router /douyin/relation/follow/list [get]
// @Summary 关注列表
// @Tags Relation
// @Accept json
// @Produce json
// @Param user_id query int64 true "用户"
// @Success      200
func GetFollowingList(ctx *gin.Context) {
	//bind arg and validate arg
	userId, err := getUserId(ctx)
	if err != nil {
		utils.LogParamError(ctx, err)
		return
	}

	followingList := service.NewFollowService(ctx).GetFollowingList(userId)
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(http.StatusOK, &FollowListResp{
		Errno:    errno.Success,
		UserList: followingList,
	})
}

func getUserId(ctx *gin.Context) (int64, error) {
	userIdS, ok := ctx.GetQuery("user_id")
	if !ok {
		return 0, errors.New("has not userId")
	}
	userId, err := strconv.ParseInt(userIdS, 10, 64)
	if err != nil {
		return 0, errors.New("userId has err format")
	}
	return userId, nil
}

// GetFollowersList
//
// @Router /douyin/relation/follower/list [get]
// @Summary 粉丝列表
// @Tags Relation
// @Accept json
// @Produce json
// @Param user_id query int64 true  "用户"
// @Success      200
func GetFollowersList(ctx *gin.Context) {
	//bind arg and validate arg
	userId, err := getUserId(ctx)
	if err != nil {
		utils.LogParamError(ctx, err)
		return
	}

	followingList := service.NewFollowService(ctx).GetFollowerList(userId)
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(http.StatusOK, &FollowListResp{
		Errno:    errno.Success,
		UserList: followingList,
	})
}

// GetFriendList
//
// @Router /douyin/relation/friend/list [get]
// @Summary 好友列表
// @Description 登录用户在消息页展示已关注的用户列表
// @Tags Message
// @Accept json
// @Produce json
// @Param user_id query int64 true  "用户"
// @Success      200
func GetFriendList(c *gin.Context) {

}
