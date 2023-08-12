package service

type FollowService interface {
	// IsFollowing 根据当前用户id和目标用户id来判断当前用户是否关注了目标用户
	IsFollowing(userId int64, targetId int64) (bool, error)
	// GetFollowerCount 根据用户id来查询用户被多少其他用户关注
	GetFollowerCount(userId int64) (int64, error)
	// GetFollowingCount 根据用户id来查询用户关注了多少其它用户
	GetFollowingCount(userId int64) (int64, error)
	// AddFollowRelation 当前用户关注目标用户
	AddFollowRelation(userId int64, targetId int64) (bool, error)
	// DeleteFollowRelation 当前用户取消对目标用户的关注
	DeleteFollowRelation(userId int64, targetId int64) (bool, error)
	// GetFollowing 获取当前用户的关注列表
	GetFollowing(userId int64) ([]UserRespond, error)
	// GetFollowers 获取当前用户的粉丝列表
	GetFollowers(userId int64) ([]UserRespond, error)
}
