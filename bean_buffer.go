package bean

import (
	"errors"
	"fmt"
	"github.com/aivyss/typex/collection"
	"reflect"
	"sync"
)

type beanBuffer struct {
	bufferOnce     sync.Once
	constructorMap map[reflect.Type]any
}

func GetBeanBuffer() *beanBuffer {
	return &beanBuffer{
		constructorMap: map[reflect.Type]any{},
	}
}

// RegisterBeanLazy
// constructor1: func(arg1, arg2, arg3, ....) (BeanType, error)
// constructor2: func(arg1, arg2, arg3, ....) BeanType
func (b *beanBuffer) RegisterBeanLazy(constructor any) {
	beanType := reflect.TypeOf(constructor).Out(0)
	b.constructorMap[beanType] = constructor
}

// RegisterBeansLazy
// constructor1: func(arg1, arg2, arg3, ....) (BeanType, error)
// constructor2: func(arg1, arg2, arg3, ....) BeanType
func (b *beanBuffer) RegisterBeansLazy(constructors ...any) {
	collection.ForEach(constructors, func(constructor any) {
		b.RegisterBeanLazy(constructor)
	})
}

func (b *beanBuffer) StartLazyLoading() error {
	var err error
	isAlreadyInitialized := true

	b.bufferOnce.Do(func() {
		dependencyTrees := b.getDependencyTrees()
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

func (b *beanBuffer) getDependencyTrees() []collection.DescendNode[reflect.Type] {
	var result []collection.DescendNode[reflect.Type]

	for t, constructor := range b.constructorMap {
		var leaves []collection.DescendNode[reflect.Type]
		typeOf := reflect.TypeOf(constructor)
		for i := 0; i < typeOf.NumIn(); i++ {
			leaf := b.recursiveDependencyTree(typeOf.In(i))
			leaves = append(leaves, leaf)
		}

		node := collection.NewDescendNode(t)
		for _, leaf := range leaves {
			node.AddDescendantNode(leaf)
		}

		result = append(result, node)
	}

	return result
}

func (b *beanBuffer) recursiveDependencyTree(in reflect.Type) collection.DescendNode[reflect.Type] {
	var trees []collection.DescendNode[reflect.Type]

	if constructor, ok := b.constructorMap[in]; ok {
		typeOf := reflect.TypeOf(constructor)
		for i := 0; i < typeOf.NumIn(); i++ {
			leaf := b.recursiveDependencyTree(typeOf.In(i))
			trees = append(trees, leaf)
		}

		node := collection.NewDescendNode(in)
		for _, leaf := range trees {
			node.AddDescendantNode(leaf)
		}

		return node
	}

	return collection.NewDescendNode(in)
}

func (b *beanBuffer) registerBeanRecursive(tree collection.DescendNode[reflect.Type]) error {
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

	if constructor, ok := b.constructorMap[tree.This()]; ok {
		return RegisterBeanEager(constructor)
	}

	return errors.New(fmt.Sprintf("fail to create bean: %s\n", tree.This().String()))
}
