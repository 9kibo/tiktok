package main

import "C"
import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"tiktok/biz/config"
	"tiktok/biz/handler"
	"tiktok/biz/middleware/ginmw"
	"time"
)

func initRouter(e *gin.Engine) {
	apiRouter := e.Group("/douyin")
	{
		withNoLogin := apiRouter.Group("", WithJwtAuth())
		withNoLogin.GET("feed", handler.Feed)
		withNoLogin.POST("user/register", handler.Register)
		withNoLogin.POST("user/login", handler.Login)
	}
	{
		apiRouter.POST("publish/action", handler.UpVideo)
		apiRouter.GET("publish/list", handler.VideoList)
		apiRouter.GET("user/", handler.UserInfo)
		//互动
		apiRouter.POST("favorite/action", handler.Favorite)
		apiRouter.GET("favorite/list", handler.FavoriteList)
		apiRouter.POST("comment/action", handler.CommAction)
		apiRouter.GET("comment/list", handler.CommList)
		//社交
		apiRouter.POST("relation/action", handler.RelationAction)
		apiRouter.GET("relation/follow/list", handler.GetFollowingList)
		apiRouter.GET("relation/follower/list", handler.GetFollowersList)
		apiRouter.GET("relation/friend/list", handler.GetFriendList)
	}
}
func WithJwtAuth() gin.HandlerFunc {
	jwtConfig := config.C.Jwt
	if SecretKeyBytes, err := base64.StdEncoding.DecodeString(config.C.Jwt.SecretKey); err != nil {
		panic("Jwt.SecretKey must a base64, err=" + err.Error())
	} else {
		return ginmw.WithJwtAuth(&ginmw.Jwt{
			Alg:       jwt.GetSigningMethod(jwtConfig.Alg),
			SecretKey: SecretKeyBytes,
			TokenKey:  jwtConfig.TokenKey,
			Issuer:    jwtConfig.Issuer,
			Audience:  jwtConfig.Audience,
			ExpireDay: time.Duration(jwtConfig.ExpireDay) * 24 * time.Hour,
		})
	}
}
