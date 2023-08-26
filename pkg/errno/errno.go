package errno

import (
	"errors"
	"fmt"
)

func NewErrno(code int, msg string) Errno {
	return Errno{code, msg}
}

// Errno 详细响应码
type Errno struct {
	Code int    `json:"status_code,omitempty" example:"1"`
	Msg  string `json:"status_msg,omitempty" example:"xxxx"`
}

func (e Errno) Error() string {
	return fmt.Sprintf("code=%d, msg=%s", e.Code, e.Msg)
}

var (
	empty = Errno{}
)

func (e Errno) IsEmpty() bool {
	return e == empty
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

	s := Service
	s.Msg = err.Error()
	return s
}
