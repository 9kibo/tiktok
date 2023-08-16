package service

type CommentService interface {
	// GetCommCountFromVId 根据视频id获取评论数量
	GetCommCountFromVId(id int64) (int64 error)
}
type CommentServiceImpl struct {
}

//全部实现后再取消, 因为返回的是FavoriteService
//func NewFavoriteService(c *gin.Context) FavoriteService {
//	return &FavoriteServiceImpl{
//		c: c,
//	}
//}
