package errno

var (
	Success    = NewErrno(SuccessCode, "Success")
	ServiceErr = NewErrno(ServiceErrCode, "Service is unable to start successfully")
	ParamErr   = NewErrno(ParamErrCode, "Wrong Parameter has been given")

	UserAlreadyExistErr             = NewErrno(UserAlreadyExistErrCode, "User already exists")
	AuthorizationFailedErr          = NewErrno(AuthorizationFailedErrCode, "Authorization failed")
	UserIsNotExistErr               = NewErrno(UserIsNotExistErrCode, "user is not exist")
	PasswordIsNotVerified           = NewErrno(AuthorizationFailedErrCode, "username or password not verified")
	FollowRelationAlreadyExistErr   = NewErrno(FollowRelationAlreadyExistErrCode, "Follow Relation already exist")
	FollowRelationNotExistErr       = NewErrno(FollowRelationNotExistErrCode, "Follow Relation does not exist")
	FavoriteRelationAlreadyExistErr = NewErrno(FavoriteRelationAlreadyExistErrCode, "Favorite Relation already exist")
	FavoriteRelationNotExistErr     = NewErrno(FavoriteRelationNotExistErrCode, "FavoriteRelationNotExistErr")
	FavoriteActionErr               = NewErrno(FavoriteActionErrCode, "favorite add failed")

	MessageAddFailedErr       = NewErrno(MessageAddFailedErrCode, "message add failed")
	FriendListNoPermissionErr = NewErrno(FriendListNoPermissionErrCode, "You can't query his friend list")
	VideoIsNotExistErr        = NewErrno(VideoIsNotExistErrCode, "video is not exist")
	CommentIsNotExistErr      = NewErrno(CommentIsNotExistErrCode, "comment is not exist")
)
