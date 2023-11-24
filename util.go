package bean

import "reflect"

func getGenericType[T any]() reflect.Type {
	return reflect.TypeOf(func(t T) {}).In(0)
}
