package validate

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"reflect"
)

// InitValidateWrapper if msgTag is null, use defaultMsgTag
func InitValidateWrapper(msgTag string) {
	if msgTag == "" {
		msgTag = defaultMsgTag
	}
	v := &defaultValidator{
		validate: &validateWrapper{
			Validate:    validator.New(),
			msgTag:      msgTag,
			msgCache:    make(map[string]string, 10),
			structCache: make(map[string]struct{}, 10),
		},
	}
	v.validate.SetTagName("binding")
	binding.Validator = v
}

// defaultValidator has not change except for remove method of lazyinit
// just to implement StructValidator to assign to binding.Validator
type defaultValidator struct {
	validate *validateWrapper
}

// ValidateStruct some with binding.defaultValidator
func (v *defaultValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(binding.SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct  some with binding.defaultValidator
func (v *defaultValidator) validateStruct(obj any) error {
	return v.validate.Struct(obj)
}

// Engine some with binding.defaultValidator
func (v *defaultValidator) Engine() any {
	return v.validate
}
