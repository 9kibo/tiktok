package model

type Video struct {
	Id        int64 `json:"id,omitempty"`
	CreatedAt int64 `json:"-"`
	UpdatedAt int64 `json:"-"`
	DeletedAt int64 `json:"-" `

	AuthorId int64  `json:"-" `
	Title    string `json:"title,omitempty"`
	PlayUrl  string `json:"play_url,omitempty"`
	CoverUrl string `json:"cover_url,omitempty" `

	FavoriteCount int32 `json:"favorite_count" gorm:"-"`
	CommentCount  int32 `json:"comment_count" gorm:"-"`

	//非表
	IsFavorite bool  `json:"is_favorite" gorm:"-"`
	Author     *User `json:"author" gorm:"-"`
}
type VideoFavor struct {
	Id        int64
	CreatedAt int64
	UserId    int64
	VideoId   int64
}
