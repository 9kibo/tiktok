package service

type CommentInfo struct {
	Id         int64       `json:"id"`
	UserInfo   UserRespond `json:"user"`
	Content    string      `json:"content"`
	CreateDate string      `json:"create_date"`
}
type CommentService interface {
	// GetCommCountFromVId 根据视频id获取评论数量
	GetCommCountFromVId(id int64) (int64 error)
}
