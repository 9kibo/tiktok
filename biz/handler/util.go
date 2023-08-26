package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/model"
)

func respOfUpdate(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, model.RespSuccess)
}
