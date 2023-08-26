package model

import (
	"regexp"
	"tiktok/pkg/errno"
)

// User
// Avatar default constant.DefaultUserAvatar
// BackgroundImage default constant.DefaultUserBackgroundImage
type User struct {
	Id        int64 `json:"id,omitempty"`
	CreatedAt int64 `json:"-"`
	DeletedAt int64 `json:"-"`

	Username        string `json:"name,omitempty"`
	Password        string `json:"-"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`

	FollowingCount int64 `json:"follow_count" gorm:"-"`
	FollowerCount  int64 `json:"follower_count" gorm:"-"`
	//获赞数量
	TotalFavorited int64 `json:"total_favorited" gorm:"-"`
	//点赞数量
	FavoriteCount int64 `json:"favorite_count" gorm:"-"`
	WorkCount     int64 `json:"work_count" gorm:"-"`

	//非表
	IsFollow bool `json:"is_follow" gorm:"-"`
}

// UserSignReq 登录或者注册请求参数
type UserSignReq struct {
	//最少1位, 必须是中文、英文、数字包括下划线
	Username string `form:"username"  binding:"gte=1,lte=32,regexp=^[\u4e00-\u9fa5\\w_]+$" errMsg:"最少1位, 必须是中文、英文、数字包括下划线"`
	//至少6位, 至少包含1个大写字母, 1个小写字母, 1个数字, 1个特殊字符
	Password string `form:"password"  binding:"gte=6,lte=32" errMsg:"至少6位"`
}

var (
	upperRegex, _   = regexp.Compile("[A-Z]")
	lowerRegex, _   = regexp.Compile("[a-z]")
	numberRegex, _  = regexp.Compile("\\d")
	specialRegex, _ = regexp.Compile("[~!@#$%^&*()_+`\\-=[\\]{}:;\"'|\\\\<,>.?/]")
)

func (r *UserSignReq) CheckPwd() error {
	if !upperRegex.MatchString(r.Password) {
		return errno.PasswordIsNotVerified.AppendMsg("at least 1 uppercase letter")
	}

	if !lowerRegex.MatchString(r.Password) {
		return errno.PasswordIsNotVerified.AppendMsg("at least 1 lowercase letter")
	}

	if !numberRegex.MatchString(r.Password) {
		return errno.PasswordIsNotVerified.AppendMsg(" at least 1 digit")
	}

	if !specialRegex.MatchString(r.Password) {
		return errno.PasswordIsNotVerified.AppendMsg("at least 1 special character in ( ~!@#$%^&*()_+`-=[]{}:;\"'|\\<,>.?/ )")
	}
	return nil
}
