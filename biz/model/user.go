package model

// User
// Avatar 默认 https://th.bing.com/th/id/OIP.FlsXXU-wKCf-Us4zDBfvzwHaHa?pid=ImgDet&rs=1
// BackgroundImage 默认 https://th.bing.com/th/id/OIP.3C4Yst9hMZpZFOh2kCzNFwAAAA?pid=ImgDet&rs=1
type User struct {
	Id        int64 `json:"id,omitempty"`
	CreatedAt int64 `json:"-"`
	DeletedAt int64 `json:"-"`

	Username        string       `json:"name,omitempty"`
	Password        string       `json:"-"`
	Avatar          string       `json:"avatar"`
	BackgroundImage string       `json:"background_image"`
	Signature       string       `json:"signature"`
	FavoriteVideos  []VideoFavor `gorm:"goreirgkey:UserId"` //一对多

	FollowCount   int64 `json:"follow_count" gorm:"-"`
	FollowerCount int64 `json:"follower_count" gorm:"-"`
	//获赞数量
	TotalFavorited int64 `json:"total_favorited" gorm:"-"`
	//点赞数量
	FavoriteCount int64 `json:"favorite_count" gorm:"-"`
	WorkCount     int64 `json:"work_count" gorm:"-"`

	//非表
	IsFollow bool `json:"is_follow" gorm:"-"`
}
