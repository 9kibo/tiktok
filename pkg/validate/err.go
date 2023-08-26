package validate

import (
	"strings"
)

type Err struct {
	Field  string `json:"field"`
	ErrMsg string `json:"errMsg"`
}

func (err Err) Error() string {
	return err.Field + ": " + err.ErrMsg
}

type Errs []Err

func (errs Errs) Error() string {
	sb := strings.Builder{}
	errLen := len(errs)
	for i, err := range errs {
		sb.WriteString(err.Field)
		sb.WriteString(": ")
		sb.WriteString(err.ErrMsg)
		if i+1 != errLen {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}
