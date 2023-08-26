package model

import (
	"tiktok/pkg/errno"
)

// BaseResp 基础响应体
type BaseResp struct {
	//业务响应码
	Code int `json:"status_code,omitempty" example:"1"`
	//业务消息
	Msg string `json:"status_msg,omitempty" example:"xxxx"`
}

var (
	RespSuccess = *baseResp(errno.Success)
)

// BuildBindResp for gin bind struct mapping request args err
func BuildBindResp(err error) *BaseResp {
	return baseResp(errno.Param.WithMessage(err.Error()))
}

// BuildBaseResp convert error and build BaseResp
func BuildBaseResp(err error) *BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}
	return baseResp(errno.ConvertErr(err))
}

// baseResp build BaseResp from error
func baseResp(err errno.Errno) *BaseResp {
	return &BaseResp{
		Code: err.Code,
		Msg:  err.Msg,
	}
}
