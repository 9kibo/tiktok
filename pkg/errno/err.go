package errno

var (
	Success   = NewErrno(SuccessCode, "Success")
	Service   = NewErrno(ServiceCode, "Service is unable to start successfully")
	Param     = NewErrno(ParamCode, "Wrong Parameter has been given")
	Update    = NewErrno(ServiceCode, "Update something but fail")
	NotExists = NewErrno(ServiceCode, "Update something but fail")

	UserAlreadyExist      = NewErrno(UserAlreadyExistCode, "User already exists")
	AuthorizationFailed   = NewErrno(AuthorizationFailedCode, "Authorization failed")
	PasswordIsNotVerified = NewErrno(AuthorizationFailedCode, "username or password not verified")

	FollowAlreadyExist = NewErrno(FollowRelationAlreadyExistErrCode, "Follow Relation already exist")
	FollowNotExist     = NewErrno(FollowRelationNotExistErrCode, "Follow Relation does not exist")

	FavoriteRelationAlreadyExistErr = NewErrno(FavoriteRelationAlreadyExistErrCode, "Favorite Relation already exist")
	FavoriteRelationNotExistErr     = NewErrno(FavoriteRelationNotExistErrCode, "FavoriteRelationNotExistErr")
	FavoriteActionErr               = NewErrno(FavoriteActionErrCode, "favorite add failed")

	MessageAddFailedErr       = NewErrno(MessageAddFailedErrCode, "message add failed")
	FriendListNoPermissionErr = NewErrno(FriendListNoPermissionErrCode, "You can't query his friend list")

	VideoIsNotExistErr = NewErrno(VideoIsNotExistErrCode, "video is not exist")

	CommentIsNotExistErr = NewErrno(CommentIsNotExistErrCode, "comment is not exist")

	MessageChatToUserNotExist = NewErrno(MessageChatToUserNotExistCode, "comment is not exist")
)
