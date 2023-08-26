package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/dao"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

type MessageService interface {
	MessageAction(req *model.MessageActionReq)
	MessageChatRecord(req *model.MessageChatRecordReq) []*model.Message
}
type MessageServiceImpl struct {
	ctx *gin.Context
}

func NewMessageService(ctx *gin.Context) MessageService {
	return &MessageServiceImpl{
		ctx: ctx,
	}
}
func (s *MessageServiceImpl) MessageAction(req *model.MessageActionReq) {
	userId := s.ctx.MustGet(constant.UserId).(int64)
	err := dao.AddMessage(&model.Message{
		Content:    req.Content,
		FromUserId: userId,
		ToUserId:   req.ToUserID,
	})
	if err != nil {
		utils.LogDB(s.ctx, err)
		return
	}
}
func (s *MessageServiceImpl) MessageChatRecord(req *model.MessageChatRecordReq) []*model.Message {
	userId := s.ctx.MustGet(constant.UserId).(int64)
	exists, err := dao.ExistsUserById(req.ToUserId)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return nil
	} else if !exists {
		utils.LogBizErr(s.ctx, errno.MessageChatToUserNotExist,
			http.StatusOK, fmt.Sprintf("user[%d] chat to not exists user[%d]", userId, req.ToUserId))
		return nil
	}
	messageList, err := dao.GetMessageRecordByLastTime(userId, req.ToUserId, req.PreMsgTime)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return nil
	}
	return messageList
}
