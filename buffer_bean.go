package bean

import (
	"reflect"
	"sort"
	"sync"
)

var bufferOnce sync.Once
var buffCreateOnce sync.Once
var buff *beanBuffer = nil

type constructorSet struct {
	constructor any
	paramNum    int
	args        []any
}

type beanBuffer struct {
	constructors []any
	args         [][]any
}

func GetBeanBuffer() *beanBuffer {
	buffCreateOnce.Do(func() {
		buff = &beanBuffer{}
	})

	return buff
}

func (b *beanBuffer) RegisterBean(constructor any) {
	b.constructors = append(b.constructors, constructor)
	b.args = append(b.args, []any{})
}

func (b *beanBuffer) RegisterBeanWithArgs(constructor any, args ...any) {
	b.constructors = append(b.constructors, constructor)
	b.args = append(b.args, args)
}

func (b *beanBuffer) Buffer() []error {
	errs := make([]error, 0, len(b.constructors))

	bufferOnce.Do(func() {
		var constructors []constructorSet
		for idx, constructor := range b.constructors {
			paramNum := reflect.TypeOf(constructor).NumIn()
			constructors = append(constructors, constructorSet{
				constructor: constructor,
				args:        b.args[idx],
				paramNum:    paramNum,
			})
		}

		sort.Slice(constructors, func(i, j int) bool {
			return constructors[i].paramNum < constructors[j].paramNum
		})

		for _, constructor := range constructors {
			err := RegisterBeanWithArgs(constructor.constructor, constructor.args)
			errs = append(errs, err)
		}
	})

	return errs
}
