package model

type Follow struct {
	Id         int64
	CreatedAt  int64
	FolloweeId int64
	FollowerId int64
}

type FollowActionReq struct {
	UserId   int64 `form:"user_id" binding:"gt=0"  errMsg:"非法数字"`
	ToUserId int64 `form:"to_user_id" binding:"gt=0"  errMsg:"非法数字"`
	// 1-关注，2-取消关注
	ActionType int32 `form:"action_type" binding:"gte=1,lte=2"  errMsg:"没有该action"`
}

func (t FollowActionReq) GetFollow() *Follow {
	return &Follow{
		FolloweeId: t.ToUserId,
		FollowerId: t.UserId,
	}
}

type FriendUser struct {
	*User
	// 和该好友的最新聊天消息
	Message string `json:"message"`
	// message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	MsgType int64 `json:"msgType"`
}
