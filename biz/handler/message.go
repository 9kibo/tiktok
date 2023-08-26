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
