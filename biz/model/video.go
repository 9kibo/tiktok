package model

import "mime/multipart"

type Video struct {
	Id        int64 `json:"id,omitempty"`
	CreatedAt int64 `json:"-"`
	UpdatedAt int64 `json:"-"`
	DeletedAt int64 `json:"-" `

	AuthorId int64  `json:"-" `
	Author   *User  `json:"author" gorm:"-"`
	Title    string `json:"title,omitempty"`
	PlayUrl  string `json:"play_url,omitempty"`
	CoverUrl string `json:"cover_url,omitempty" `

	FavoriteCount int32 `json:"favorite_count" gorm:"-"`
	CommentCount  int32 `json:"comment_count" gorm:"-"`

	//非表
	IsFavorite bool `json:"is_favorite" gorm:"-"`
}

type VideoUploadReq struct {
	Data  multipart.File
	File  *multipart.FileHeader
	Token string `form:"token" binding:"required"`
	Title string `form:"title" binding:"required"`
}
