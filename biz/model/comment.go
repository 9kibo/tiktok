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
type CommInfo struct {
	UserId  int64  `redis:"userid"`
	Content string `redis:"content"`
}

type CommReq struct {
	UserId    int64
	VideoId   int64  `form:"video_id" binging:"gt=0"`
	Action    int64  `form:"action_type" binging:"oneof=1 2"`
	Text      string `form:"comment_text" binging:"omitempty"`
	DelCommId int64  `form:"comment_id" binging:"omitempty"`
}
type CommentsReq struct {
	UserId  int64
	VideoId int64 `form:"video_id" binging:"gt=0"`
}

var ErrComm = Comment{
	Id:        0,
	CreatedAt: 0,
	Content:   "错误信息",
	User:      nil,
}
