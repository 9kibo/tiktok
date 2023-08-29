package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/pkg/constant"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

// isTokenUser userId是否和token中的一样, 不一样会自动响应
func isTokenUser(ctx *gin.Context, userId int64) bool {
	//检查是否和token中的一样
	uid := ctx.MustGet(constant.UserId).(int64)
	if uid != userId {
		utils.LogBizErr(ctx, errno.AuthorizationFailed, http.StatusConflict, "userId is not equal with uid from token, maybe a attack")
		return false
	}
	return true
}
