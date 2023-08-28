package errno

const (
	SuccessCode = 0
	ServiceCode = iota + 1000
	ParamCode
	AuthorizationFailedCode

	UserAlreadyExistCode

	FollowRelationAlreadyExistErrCode
	FollowRelationNotExistErrCode

	FavoriteRelationAlreadyExistErrCode
	FavoriteRelationNotExistErrCode
	FavoriteActionErrCode

	MessageAddFailedErrCode
	FriendListNoPermissionErrCode

	VideoIsNotExistErrCode
	CommentIsNotExistErrCode
	CommentActionErrCode

	MessageChatToUserNotExistCode
)
