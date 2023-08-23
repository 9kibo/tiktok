package validate

import (
	"strings"
)

type ValidateErr struct {
	Field  string `json:"field"`
	ErrMsg string `json:"errMsg"`
}

func (err ValidateErr) Error() string {
	return err.Field + ": " + err.ErrMsg
}

type ValidateErrs []ValidateErr

func (errs ValidateErrs) Error() string {
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
