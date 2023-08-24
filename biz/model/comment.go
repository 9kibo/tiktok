package model

type Comment struct {
	Id        int64 `json:"id,omitempty"`
	CreatedAt int64 `json:"create_date"`

	UserId  int64  `json:"-"`
	VideoId int64  `json:"-"`
	Content string `json:"content,omitempty"`
	User    *User  `json:"user,omitempty" gorm:"-"`
}
type CommToJson struct {
	CommId   int64
	CreateAt int64
	UserId   int64
	VideoId  int64
	Content  string
}

type CommReq struct {
	UserId    int64
	VideoId   int64  `from:"video_id" binging:"gt=0"`
	Action    int64  `from:"action_type" binging:"oneof=1 2"`
	Text      string `from:"comment_text" binging:"omitempty"`
	DelCommId int64  `from:"comment_id" binging:"omitempty"`
}
type CommentsReq struct {
	UserId  int64
	VideoId int64 `from:"video_id" binging:"gt=0"`
}
