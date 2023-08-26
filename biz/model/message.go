package model

type Message struct {
	Id        int64 `json:"id"`
	CreatedAt int64 `json:"create_time"`
	// 该消息接收者的id
	ToUserId int64 `json:"to_user_id"`
	// 该消息发送者的id
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"content"`
}

type MessageActionReq struct {
	// 1-发送消息
	ActionType int32 `form:"action_type" binding:"eq=1" errMsg:"没有该action"`
	// 消息内容
	Content string `form:"content"  binding:"ne=0" errMsg:"必须有消息内容"`
	// 对方用户id
	ToUserID int64 `form:"to_user_id" binding:"gt=0" errMsg:"非法数字"`
}
type MessageChatRecordReq struct {
	// 对方用户id
	ToUserId int64 `form:"to_user_id"  binding:"gt=0" errMsg:"非法数字"`
	//上次最新消息的时间, to find messages which is after this time
	PreMsgTime int64 `form:"pre_msg_time"  binding:"gt=0" errMsg:"非法日期"`
}
