package bean

import (
	"errors"
	"fmt"
	"github.com/aivyss/typex/collection"
	"github.com/aivyss/typex/utilx"
	"reflect"
)

var m = map[reflect.Type]any{}

func RegisterBeanEager(constructor any) error {
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

	if len(collection.Filter(beans, func(param any) bool { return param == nil })) > 0 {
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
		if r.Type() == utilx.GetGenericType[error]() && !utilx.IsNil(r.Interface()) {
			e = r.Interface().(error)
			break
		}
	}

	if e != nil {
		return e
	}

	// check bean
	isDetectedTheBean := false
	for _, r := range returns {
		if r.Type() != utilx.GetGenericType[error]() {
			m[r.Type()] = r.Interface()
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
	genericType := utilx.GetGenericType[T]()
	if b, ok := m[genericType]; ok {
		if b == nil {
			var t T
			return t, nil
		}

		return b.(T), nil
	}

	var t T
	return t, errors.New(fmt.Sprintf("no bean: %s", genericType.String()))
}

func MustGetBean[T any]() T {
	bean, err := GetBean[T]()
	if err != nil {
		panic(err)
	}

	return bean
}

func Clean() {
	m = map[reflect.Type]any{}
}
