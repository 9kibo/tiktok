package server

type FavoriteServer interface {
	//IsFavorite 根据id返回是否点赞了该视频
	IsFavorite(videoId int64, userId int64) (bool, error)
	//FavouriteCount 根据当前视频id获取当前视频点赞数量。
	FavouriteCount(videoId int64) (int64, error)
	//TotalFavourite 根据userId获取这个用户总共被点赞数量
	TotalFavourite(userId int64) (int64, error)
	//FavouriteVideoCount 根据userId获取这个用户点赞视频数量
	FavouriteVideoCount(userId int64) (int64, error)

	//当前操作行为，1点赞，2取消点赞。
	FavouriteAction(userId int64, videoId int64, actionType int32) error
	// GetFavouriteList 获取当前用户的所有点赞视频
	GetFavouriteList(userId int64, curId int64) ([]VideoRespond, error)
}
