package bean

import (
	"github.com/aivyss/typex"
	"reflect"
	"sync"
)

var bufferOnce sync.Once
var buffCreateOnce sync.Once
var buff *beanBuffer = nil

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

func (b *beanBuffer) Buffer() error {
	var err error
	bufferOnce.Do(func() {
		idxMap := map[reflect.Type]int{}
		for idx, constructor := range b.constructors {
			out := reflect.TypeOf(constructor).Out(0)
			idxMap[out] = idx
		}

		dependencyTrees := b.GetDependencyTrees()
		for _, tree := range dependencyTrees {
			if e := b.registerBeanRecursive(tree, idxMap); e != nil {
				err = e
				break
			}
		}
	})

	return err
}

func (b *beanBuffer) GetDependencyTrees() []typex.DescendNode[reflect.Type] {
	var result []typex.DescendNode[reflect.Type]
	for _, constructor := range b.constructors {
		var leaves []typex.DescendNode[reflect.Type]
		typeOf := reflect.TypeOf(constructor)
		for i := 0; i < typeOf.NumIn(); i++ {
			leaf := b.recursiveDependencyTree(typeOf.In(i))
			leaves = append(leaves, leaf)
		}

		node := typex.NewDescendNode(typeOf.Out(0))
		for _, leaf := range leaves {
			node.AddDescendantNode(leaf)
		}

		result = append(result, node)
	}

	return result
}

func (b *beanBuffer) recursiveDependencyTree(in reflect.Type) typex.DescendNode[reflect.Type] {
	var trees []typex.DescendNode[reflect.Type]
	for _, constructor := range b.constructors {
		typeOf := reflect.TypeOf(constructor)
		if typeOf.Out(0) == in {
			for i := 0; i < typeOf.NumIn(); i++ {
				leaf := b.recursiveDependencyTree(typeOf.In(i))
				trees = append(trees, leaf)
			}
		}
	}

	node := typex.NewDescendNode(in)
	for _, leaf := range trees {
		node.AddDescendantNode(leaf)
	}
	return node
}

func (b *beanBuffer) registerBeanRecursive(tree typex.DescendNode[reflect.Type], idxMap map[reflect.Type]int) error {
	if _, ok := m[tree.This()]; ok {
		return nil
	}
	if tree.HasDescendants() {
		for _, child := range tree.GetDescendants() {
			err := b.registerBeanRecursive(child, idxMap)
			if err != nil {
				return err
			}
		}
	}

	if idx, ok := idxMap[tree.This()]; ok {
		return RegisterBeanWithArgs(b.constructors[idx], b.args[idx])
	}

	return nil
}
