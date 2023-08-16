package model

type Comment struct {
	Id        int64 `json:"id,omitempty"`
	CreatedAt int64 `json:"create_date"`

	UserId  int64  `json:"-"`
	VideoId int64  `json:"-"`
	Content string `json:"content,omitempty"`

	//非表
	User *User `json:"user,omitempty" gorm:"-"`
}
