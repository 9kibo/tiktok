package model

type VideoFavor struct {
	Id        int64 `gorm:"id"`
	CreatedAt int64 `gorm:"create_at"`
	UserId    int64 `gorm:"user_id"`
	VideoId   int64 `gorm:"video_id"`
	VideoInfo Video `gorm:"foreignkey:Id"` //一对一
}
type FavoriteReq struct {
	UserId     int64 `binding:"ne=0"`
	VideoId    int64 `from:"video_id" binding:"ne=0"`
	ActionType int32 `from:"action_type" binding:"gte=1,lte=2"`
}
type FavoriteListReq struct {
	CurUserId int64 `binding:"ne=0"`
	UserId    int64 `from:"user_id" binding:"ne=0"`
}
