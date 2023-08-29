package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/model"
	"tiktok/biz/service"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

// MessageAction 登录用户在消息页展示已关注的用户列表，点击用户头像进入聊天页后可以发送消息
//
// @Router /douyin/message/action/ [post]
// @Summary 发送消息给某个关注的人
// @Schemes http
// @Tags Relation
// @Produce      json
// @Param action_type query int32 true "目前只有1-send msg"
// @Param to_user_id query int64 true "接收者"
// @Param content query string true "内容"
// @Success      200
func MessageAction(ctx *gin.Context) {
	req := model.MessageActionReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.LogParamError(ctx, err)
		return
	}
	service.NewMessageService(ctx).MessageAction(&req)
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(http.StatusOK, errno.Success)
}

type MessageChatResp struct {
	errno.Errno
	MessageList []*model.Message `json:"message_list"`
}

// MessageChatRecord chat record
//
// @Router /douyin/message/chat/ [get]
// @Summary 与某个好友的聊天记录
// @Schemes http
// @Tags Relation
// @Produce      json
// @Param pre_msg_time query int64 true "前端最新时间, 后端需要返回该时间之后的消息"
// @Param to_user_id query int64 true "好友"
// @Success      200
func MessageChatRecord(ctx *gin.Context) {
	req := model.MessageChatRecordReq{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.LogParamError(ctx, err)
		return
	}
	messageList := service.NewMessageService(ctx).MessageChatRecord(&req)
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(http.StatusOK, &MessageChatResp{
		Errno:       errno.Success,
		MessageList: messageList,
	})
}

// FriendListResp 好友列表响应
type FriendListResp struct {
	errno.Errno
	//好友列表
	UserList []*model.FriendUser `json:"user_list"`
}

// GetFriendList
//
// @Router /douyin/relation/friend/list/ [get]
// @Summary 好友列表
// @Description 用户发送过消息的人
// @Tags Relation
// @Produce json
// @Param user_id query int64 true  "用户"
// @Success      200
func GetFriendList(ctx *gin.Context) {
	//bind arg and validate arg
	userId, err := getUserId(ctx)
	if err != nil {
		utils.LogParamError(ctx, err)
		return
	}

	friendList := service.NewFollowService(ctx).GetFriendList(userId)
	if ctx.IsAborted() {
		return
	}
	ctx.JSON(http.StatusOK, &FriendListResp{
		Errno:    errno.Success,
		UserList: friendList,
	})
}
