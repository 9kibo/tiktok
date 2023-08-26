package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/dao"
	"tiktok/biz/middleware/ginmw"
	"tiktok/biz/middleware/redis"
	"tiktok/biz/model"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

type UserService interface {
	//Register 用户注册
	Register(req *model.UserSignReq) (int64, string)
	//Login 用户登录
	Login(req *model.UserSignReq) (int64, string)

	//GetUserByUserId 根据id获取user
	GetUserByUserId(userId int64) *model.User
}

func NewUserService(ctx *gin.Context) UserService {
	return &UserServiceImpl{
		ctx: ctx,
	}
}

type UserServiceImpl struct {
	ctx *gin.Context
}

func (s UserServiceImpl) Register(req *model.UserSignReq) (int64, string) {
	//check username if exists
	exists, err := dao.ExistsUserByUsername(req.Username)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return 0, ""
	} else if exists {
		utils.LogBizErr(s.ctx, errno.UserAlreadyExist, http.StatusOK, "username conflict")
		return 0, ""
	}

	//add user
	userId, err := dao.AddUser(&model.User{
		Username:        req.Username,
		Password:        req.Password,
		Avatar:          constant.DefaultUserAvatar,
		BackgroundImage: constant.DefaultUserBackgroundImage,
	})
	if err != nil {
		utils.LogDB(s.ctx, err)
		return 0, ""
	}

	return s.createToken(userId)
}

// createToken create token
func (s UserServiceImpl) createToken(userId int64) (int64, string) {
	token, err := ginmw.JWT.CreateToken(&ginmw.PublicClaims{UserId: userId})
	if err != nil {
		utils.LogWithRequestId(s.ctx, constant.LMJwt, err).Debug("create token fail")
		return 0, ""
	}

	return userId, token
}
func (s UserServiceImpl) Login(req *model.UserSignReq) (int64, string) {
	//if username and password equal
	user, err := dao.MustGetUserByUsernamePassword(req.Username, req.Password)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return 0, ""
	}

	return s.createToken(user.Id)
}

func (s UserServiceImpl) GetUserByUserId(userId int64) *model.User {
	user, err := dao.MustGetUserById(userId)
	if err != nil {
		utils.LogDB(s.ctx, err)
		return nil
	}
	followRedisService := redis.NewFollowService(s.ctx)
	var ok bool
	if user.FollowingCount, ok = followRedisService.GetFollowingCount(userId); !ok {
		user.FollowingCount, err = dao.GetFollowingCount(userId)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return nil
		}
	}
	if user.FollowerCount, ok = followRedisService.GetFollowerCount(userId); !ok {
		user.FollowerCount, err = dao.GetFollowerCount(userId)
		if err != nil {
			utils.LogDB(s.ctx, err)
			return nil
		}
	}
	var isFollow bool
	userIdS, _ := s.ctx.Get(constant.UserId)
	loginUserId := userIdS.(int64)
	if isFollow, ok = followRedisService.ExistsFollow(loginUserId, userId); !ok {
		isFollow, err = dao.ExistsFollow(&model.Follow{
			FollowerId: loginUserId,
			FolloweeId: userId,
		})
		if err != nil {
			utils.LogDB(s.ctx, err)
			return nil
		}
	}
	user.IsFollow = isFollow

	//if user.FavoriteCount, ok = followRedisService.GetFollowerCount(userId); !ok {
	//	user.FavoriteCount, err = dao.GetFollowingCount(userId)
	//}
	//if user.TotalFavorited, ok = followRedisService.GetFollowerCount(userId); !ok {
	//	user.TotalFavorited, err = dao.GetFollowingCount(userId)
	//}
	//if user.WorkCount, ok = followRedisService.GetFollowerCount(userId); !ok {
	//	user.WorkCount, err = dao.GetFollowingCount(userId)
	//}
	return user
}
