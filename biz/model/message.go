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
