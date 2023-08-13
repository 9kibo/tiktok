package errmsg

const (
	SUCCSE = 0
	ERROR  = 1
	// 1xxx 用户模块错误

	// 2xxx 视频模块错误

	// 3xxx 互动模块错误

	// 4xxx 社交模块错误

)

var statusMsg = map[int]string{
	SUCCSE: "OK",
	ERROR:  "ERR",
	// statusCode : statusMsg
}

func GetErrMsg(code int) string {
	return statusMsg[code]
}
