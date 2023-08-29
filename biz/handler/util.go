package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tiktok/biz/model"
)

func respOfUpdate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, model.RespSuccess)
}
func getUserId(ctx *gin.Context) (int64, error) {
	userIdS, ok := ctx.GetQuery("user_id")
	if !ok {
		return 0, errors.New("has not userId")
	}
	userId, err := strconv.ParseInt(userIdS, 10, 64)
	if err != nil || userId <= 0 {
		return 0, errors.New("userId has err format or value")
	}
	return userId, nil
}
