package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/model"
	"tiktok/biz/service"
	"tiktok/pkg/errno"
)

// RelationAction
//
// @Router /douyin/relation/action [post]
// @Summary 用户关系操作
// @Schemes
// @Description 关注或者取消关注用户
// @Tags Relation
// @Accept       json
// @Produce      json
// @Param user_id query int64 true "关注者"
// @Param to_user_id query int64 true "被关注者"
// @Success      200
// @Failure      400  {object}  model.BaseResp
// @Failure      403  {object}  model.BaseResp
// @Failure      500  {object}  model.BaseResp
func RelationAction(ctx *gin.Context) {
	//绑定参数并校验
	req := &model.FollowActionReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		ctx.JSON(http.StatusBadRequest, model.BuildBindResp(err))
	}

	//执行业务
	service.NewFollowService(ctx).FollowAction(req)

	//业务有错误, 不响应
	if ctx.IsAborted() {
		return
	}

	//业务成功, 响应结果
	ctx.JSON(http.StatusOK, errno.Success)
}

type FollowingListResp struct {
	model.BaseResp
	UserList []*model.User `json:"user_list"`
}

func GetFollowingList(c *gin.Context) {
}
func GetFollowersList(c *gin.Context) {
}

// GetFriendList 登录用户在消息页展示已关注的用户列表
func GetFriendList(c *gin.Context) {}
