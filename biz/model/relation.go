package model

type Follow struct {
	Id         int64
	CreatedAt  int64
	FolloweeId int64
	FollowerId int64
}

type FollowActionReq struct {
	UserId   int64 `form:"user_id" binding:"ne=0"`
	ToUserId int64 `form:"to_user_id" binding:"ne=0"`
	// 1-关注，2-取消关注
	ActionType int32 `form:"action_type" binding:"gte=1,lte=2"`
}

func (t FollowActionReq) GetFollow() *Follow {
	return &Follow{
		FolloweeId: t.UserId,
		FollowerId: t.ToUserId,
	}
}

type FollowingListReq struct {
	UserId int64 `form:"user_id" binding:"ne=0"`
}
type FriendUser struct {
	User
	// 和该好友的最新聊天消息
	Message string `json:"message"`
	// message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	MsgType int64 `json:"msgType"`
}
