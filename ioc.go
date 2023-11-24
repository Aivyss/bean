package bean

import (
	"errors"
	"github.com/aivyss/typex/util"
	"reflect"
)

var m = map[reflect.Type]any{}

func RegisterBeanWithArgs[T any](constructor any, args ...any) error {
	fnTypeOf := reflect.TypeOf(constructor)

	paramTypes := make([]reflect.Type, 0, fnTypeOf.NumIn())
	for i := 0; i < fnTypeOf.NumIn(); i++ {
		paramTypes = append(paramTypes, fnTypeOf.In(i))
	}

	beans := make([]any, fnTypeOf.NumIn(), fnTypeOf.NumIn())
	var nilIdxs []int
	for i, paramType := range paramTypes {
		if b, ok := m[paramType]; ok {
			beans[i] = b
		} else {
			beans[i] = nil
			nilIdxs = append(nilIdxs, i)
		}
	}

	if len(nilIdxs) > 0 {
		idx := 0
		for _, nilIdx := range nilIdxs {
			beans[nilIdx] = args[idx]
			idx += 1
		}
	}

	if len(util.Filter(beans, func(param any) bool { return param == nil })) > 0 {
		return errors.New("no bean")
	}

	params := make([]reflect.Value, fnTypeOf.NumIn(), fnTypeOf.NumIn())
	for i, b := range beans {
		params[i] = reflect.ValueOf(b)
	}

	returns := reflect.ValueOf(constructor).Call(params)

	// check error
	var e error
	errorTypeOf := reflect.TypeOf(e)
	isError := false
	for _, r := range returns {
		if r.Type() == errorTypeOf && !util.IsNil(r.Interface()) {
			e = r.Interface().(error)
			isError = true
			break
		}
	}

	if isError {
		return e
	}

	// check bean
	beanTargetTypeOf := getGenericType[T]()

	isDetectedTheBean := false
	for _, r := range returns {
		if r.Type() == beanTargetTypeOf {
			m[beanTargetTypeOf] = r.Interface()
			isDetectedTheBean = true
			break
		}
	}

	if !isDetectedTheBean {
		return errors.New("no bean is detected")
	}

	return nil
}

func RegisterBean[T any](constructor any) error {
	fnTypeOf := reflect.TypeOf(constructor)

	paramTypes := make([]reflect.Type, 0, fnTypeOf.NumIn())
	for i := 0; i < fnTypeOf.NumIn(); i++ {
		paramTypes = append(paramTypes, fnTypeOf.In(i))
	}

	beans := make([]any, fnTypeOf.NumIn(), fnTypeOf.NumIn())
	for i, paramType := range paramTypes {
		if b, ok := m[paramType]; ok {
			beans[i] = b
		} else {
			beans[i] = nil
		}
	}

	if len(util.Filter(beans, func(param any) bool { return param == nil })) > 0 {
		return errors.New("no bean")
	}

	params := make([]reflect.Value, fnTypeOf.NumIn(), fnTypeOf.NumIn())
	for i, b := range beans {
		params[i] = reflect.ValueOf(b)
	}

	returns := reflect.ValueOf(constructor).Call(params)

	// check error
	var e error = nil
	for _, r := range returns {
		if r.Type() == getGenericType[error]() && !util.IsNil(r.Interface()) {
			e = r.Interface().(error)
			break
		}
	}

	if e != nil {
		return e
	}

	// check bean
	beanTargetTypeOf := getGenericType[T]()

	isDetectedTheBean := false
	for _, r := range returns {
		if r.Type() == beanTargetTypeOf {
			m[beanTargetTypeOf] = r.Interface()
			isDetectedTheBean = true
			break
		}
	}

	if !isDetectedTheBean {
		return errors.New("no bean is detected")
	}

	return nil
}

func GetBean[T any]() (T, error) {
	genericType := getGenericType[T]()
	if b, ok := m[genericType]; ok {
		return b.(T), nil
	}

	var t T
	return t, errors.New("no bean")
}

func Clean() {
	m = map[reflect.Type]any{}
}
