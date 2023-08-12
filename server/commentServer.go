package server

type CommentServer interface {
	// GetCommCountFromVId 根据视频id获取评论数量
	GetCommCountFromVId(id int64) (int64 error)
}
