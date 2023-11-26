package bean

import (
	"errors"
	"fmt"
	"github.com/aivyss/typex"
	"reflect"
	"sync"
)

var buffCreateOnce sync.Once
var buff *beanBuffer = nil

type beanBuffer struct {
	bufferOnce     sync.Once
	constructorMap map[reflect.Type]struct {
		constructor any
		args        []any
	}
}

func GetBeanBuffer() *beanBuffer {
	buffCreateOnce.Do(func() {
		buff = &beanBuffer{
			constructorMap: map[reflect.Type]struct {
				constructor any
				args        []any
			}{},
		}
	})

	return buff
}

func (b *beanBuffer) RegisterBean(constructor any) {
	beanType := reflect.TypeOf(constructor).Out(0)
	b.constructorMap[beanType] = struct {
		constructor any
		args        []any
	}{constructor: constructor, args: []any{}}
}

func (b *beanBuffer) RegisterBeanWithArgs(constructor any, args ...any) {
	beanType := reflect.TypeOf(constructor).Out(0)
	b.constructorMap[beanType] = struct {
		constructor any
		args        []any
	}{constructor: constructor, args: args}
}

func (b *beanBuffer) Buffer() error {
	var err error
	isAlreadyInitialized := true

	b.bufferOnce.Do(func() {
		dependencyTrees := b.GetDependencyTrees()
		for _, tree := range dependencyTrees {
			if e := b.registerBeanRecursive(tree); e != nil {
				err = e
				break
			}
		}

		// release constructor resource
		isAlreadyInitialized = false
		b.constructorMap = nil
	})

	if isAlreadyInitialized {
		err = errors.New("beans already is registered")
	}

	return err
}

func (b *beanBuffer) GetDependencyTrees() []typex.DescendNode[reflect.Type] {
	var result []typex.DescendNode[reflect.Type]

	for t, constructorInfo := range b.constructorMap {
		var leaves []typex.DescendNode[reflect.Type]
		typeOf := reflect.TypeOf(constructorInfo.constructor)
		for i := 0; i < typeOf.NumIn(); i++ {
			leaf := b.recursiveDependencyTree(typeOf.In(i))
			leaves = append(leaves, leaf)
		}

		node := typex.NewDescendNode(t)
		for _, leaf := range leaves {
			node.AddDescendantNode(leaf)
		}

		result = append(result, node)
	}

	return result
}

func (b *beanBuffer) recursiveDependencyTree(in reflect.Type) typex.DescendNode[reflect.Type] {
	var trees []typex.DescendNode[reflect.Type]

	if constructorInfo, ok := b.constructorMap[in]; ok {
		typeOf := reflect.TypeOf(constructorInfo.constructor)
		for i := 0; i < typeOf.NumIn(); i++ {
			leaf := b.recursiveDependencyTree(typeOf.In(i))
			trees = append(trees, leaf)
		}

		node := typex.NewDescendNode(in)
		for _, leaf := range trees {
			node.AddDescendantNode(leaf)
		}

		return node
	}

	return typex.NewDescendNode(in)
}

func (b *beanBuffer) registerBeanRecursive(tree typex.DescendNode[reflect.Type]) error {
	if _, ok := m[tree.This()]; ok {
		return nil
	}

	if tree.HasDescendants() {
		for _, child := range tree.GetDescendants() {
			err := b.registerBeanRecursive(child)
			if err != nil {
				return err
			}
		}
	}

	if constructorInfo, ok := b.constructorMap[tree.This()]; ok {
		return RegisterBeanWithArgs(constructorInfo.constructor, constructorInfo.args)
	}

	return errors.New(fmt.Sprintf("fail to create bean: %s\n", tree.This().String()))
}
