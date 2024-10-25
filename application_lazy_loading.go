package bean

import (
	"errors"
	"github.com/aivyss/typex/collection"
	"reflect"
)

var bufferBeanApplication = GetBeanBuffer()

// RegisterBeansLazy
// constructor1: func(arg1, arg2, arg3, ....) (BeanType, error)
// constructor2: func(arg1, arg2, arg3, ....) BeanType
func RegisterBeansLazy(constructors ...any) {
	collection.ForEach(constructors, func(constructor any) {
		bufferBeanApplication.RegisterBeanLazy(constructor)
	})
}

// RegisterBeanLazy
// constructor1: func(arg1, arg2, arg3, ....) (BeanType, error)
// constructor2: func(arg1, arg2, arg3, ....) BeanType
func RegisterBeanLazy(constructor any) {
	beanType := reflect.TypeOf(constructor).Out(0)
	bufferBeanApplication.constructorMap[beanType] = constructor
}

func StartLazyLoading() error {
	var err error
	isAlreadyInitialized := true

	bufferBeanApplication.bufferOnce.Do(func() {
		dependencyTrees := bufferBeanApplication.getDependencyTrees()
		for _, tree := range dependencyTrees {
			if e := bufferBeanApplication.registerBeanRecursive(tree); e != nil {
				err = e
				break
			}
		}

		// release constructor resource
		isAlreadyInitialized = false
		bufferBeanApplication.constructorMap = nil
	})

	if isAlreadyInitialized {
		err = errors.New("beans already is registered")
	}

	return err
}
