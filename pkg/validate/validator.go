package validate

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"tiktok/pkg/constant"
	"tiktok/pkg/utils"
)

func initValidator(v *validateWrapper) {
	v.SetTagName("binding")
	err := v.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		matched, err := regexp.MatchString(fl.Param(), fl.Field().String())
		if err != nil {
			utils.Log(constant.LMValidator).WithError(err).Debug("regexp Validation err")
			return false
		}
		return matched
	})
	if err != nil {
		panic(err)
	}
}
