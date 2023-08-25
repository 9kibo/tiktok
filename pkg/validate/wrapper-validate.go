package validate

import (
	"context"
	"github.com/go-playground/validator/v10"
	"reflect"
	"sync"
)

const (
	defaultMsgTag = "errMsg"
)

type validateWrapper struct {
	*validator.Validate
	msgCache    map[string]string
	structCache map[string]struct{}
	lock        sync.Mutex
	msgTag      string
}

// Struct s must a struct for binding.defaultValidator use for struct
func (v *validateWrapper) Struct(s interface{}) error {
	err := v.StructCtx(context.Background(), s)
	if err == nil {
		return nil
	}
	//if is InvalidValidationError, mean val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType)
	if ive, ok := err.(*validator.InvalidValidationError); ok {
		return ive
	}

	typeNamespace := v.cacheStructFieldMsg(s)
	errs := Errs{}
	//the validator.FieldError only is field err, has not it's struct info, so must use by the param s
	for _, err := range err.(validator.ValidationErrors) {
		msg := v.msgCache[typeNamespace+"."+err.Field()]
		if msg == "" {
			msg = err.Error()
		}
		errs = append(errs, Err{
			Field:  err.StructField(),
			ErrMsg: msg,
		})
	}
	return errs
}

// cacheStructFieldMsg use lock and double check, to cache a struct's all field'msg
func (v *validateWrapper) cacheStructFieldMsg(s any) string {
	typeOf := reflect.TypeOf(s)
	typeNamespace := typeOf.PkgPath() + "." + typeOf.Name()

	if v.structCache[typeNamespace] != struct{}{} {
		return typeNamespace
	}

	v.lock.Lock()
	defer v.lock.Unlock()
	if v.structCache[typeNamespace] != struct{}{} {
		return typeNamespace
	}

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		v.msgCache[typeNamespace+"."+field.Name] = field.Tag.Get(v.msgTag)
	}
	v.structCache[typeNamespace] = struct{}{}

	return typeNamespace
}
