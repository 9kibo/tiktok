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
		apiRouter.GET("feed", handler.Feed)
		apiRouter.POST("user/register", handler.Register)
		apiRouter.POST("user/login", handler.Login)
	}
	{
		withLoginRoute := apiRouter.Group("", WithJwtAuth())
		withLoginRoute.POST("publish/action", handler.UpVideo)
		withLoginRoute.GET("publish/list", handler.VideoList)
		withLoginRoute.GET("user/", handler.UserInfo)
		//互动
		withLoginRoute.POST("favorite/action", handler.Favorite)
		withLoginRoute.GET("favorite/list", handler.FavoriteList)
		withLoginRoute.POST("comment/action", handler.CommAction)
		withLoginRoute.GET("comment/list", handler.CommList)
		//社交
		withLoginRoute.POST("relation/action", handler.RelationAction)
		withLoginRoute.GET("relation/follow/list", handler.GetFollowingList)
		withLoginRoute.GET("relation/follower/list", handler.GetFollowersList)
		withLoginRoute.GET("relation/friend/list", handler.GetFriendList)
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
