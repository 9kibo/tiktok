package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tiktok/biz/model"
	"tiktok/biz/service"
	"tiktok/pkg/errno"
	"tiktok/pkg/utils"
)

// UserSignResp register or login resp
type UserSignResp struct {
	errno.Errno
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

// Register
//
// @Router /douyin/user/register [post]
// @Summary 注册
// @Schemes http
// @Tags User
// @Accept       json
// @Produce      json
// @Param userSignReq query model.UserSignReq true "req"
// @Success      200
// @Failure      200  {body}  errno.Errno "参数不正确"
// @Failure      401  {body}  errno.Errno "未登录"
// @Failure      200  {body}  errno.Errno "系统错误"
func Register(ctx *gin.Context) {
	//bind arg and validate arg
	req := &model.UserSignReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		utils.LogParamError(ctx, err)
		return
	}
	err := req.CheckPwd()
	if err != nil {
		utils.LogParamError(ctx, err)
		return
	}

	userId, token := service.NewUserService(ctx).Register(req)

	if ctx.IsAborted() {
		return
	}

	ctx.JSON(http.StatusOK, &UserSignResp{
		Errno:  errno.Success,
		UserId: userId,
		Token:  token,
	})
}

// Login
//
// @Router /douyin/user/login [post]
// @Summary 登录
// @Schemes http
// @Tags User
// @Accept       json
// @Produce      json
// @Param userSignReq query model.UserSignReq true "req"
// @Success      200
// @Failure      200  {body}  errno.Errno "参数不正确"
// @Failure      401  {body}  errno.Errno "未登录"
// @Failure      200  {body}  errno.Errno "系统错误"
func Login(ctx *gin.Context) {
	//bind arg and validate arg
	req := &model.UserSignReq{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		utils.LogParamError(ctx, err)
		return
	}
	err := req.CheckPwd()
	if err != nil {
		utils.LogParamError(ctx, err)
		return
	}

	userId, token := service.NewUserService(ctx).Login(req)

	if ctx.IsAborted() {
		return
	}

	ctx.JSON(http.StatusOK, &UserSignResp{
		Errno:  errno.Success,
		UserId: userId,
		Token:  token,
	})
}

// UserInfo 个人主页：支持查看用户基本信息和投稿列表，注册用户流程简化
func UserInfo(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		utils.LogParamError(ctx, err)
		return
	}

	user := service.NewUserService(ctx).GetUserByUserId(userId)
	if ctx.IsAborted() {
		return
	}

	ctx.JSON(http.StatusOK, user)
}
