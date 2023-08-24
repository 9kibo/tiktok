package errno

import (
	"errors"
	"fmt"
)

type Errno struct {
	Code int32  `json:"status_code"`
	Msg  string `json:"status_msg"`
}

func (e Errno) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}

func NewErrno(code int32, msg string) Errno {
	return Errno{code, msg}
}

// WithMessage Errno的替换msg
func (e Errno) WithMessage(msg string) Errno {
	e.Msg = msg
	return e
}

// AppendMsg 以,连接加入的msg,返回新Errno对象
func (e Errno) AppendMsg(msg string) Errno {
	return NewErrno(e.Code, e.Msg+", "+msg)
}

// ConvertErr convert error to Errno(if the error contain Errno)
func ConvertErr(err error) Errno {
	Err := Errno{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.Msg = err.Error()
	return s
}
