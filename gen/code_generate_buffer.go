package bean

import (
	"errors"
	"fmt"
	"github.com/aivyss/typex/collection"
	"github.com/dave/jennifer/jen"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

type beanBuffer struct {
	bufferOnce      sync.Once
	constructorMap  map[reflect.Type]any
	beanChecker     map[reflect.Type]bool
	initializeOrder []reflect.Type
	beanSetName     string
	filePath        string
}

func GetBeanBuffer(beanSetName string, filePath string) *beanBuffer {
	return &beanBuffer{
		constructorMap:  map[reflect.Type]any{},
		beanChecker:     map[reflect.Type]bool{},
		initializeOrder: []reflect.Type{},
		beanSetName:     beanSetName,
		filePath:        filePath,
	}
}

// RegisterBean
// constructor1: func(arg1, arg2, arg3, ....) (BeanType, error)
// constructor2: func(arg1, arg2, arg3, ....) BeanType
func (b *beanBuffer) RegisterBean(constructor any) {
	beanType := reflect.TypeOf(constructor).Out(0)
	b.constructorMap[beanType] = constructor
}

// RegisterBeans
// constructor1: func(arg1, arg2, arg3, ....) (BeanType, error)
// constructor2: func(arg1, arg2, arg3, ....) BeanType
func (b *beanBuffer) RegisterBeans(constructors ...any) {
	collection.ForEach(constructors, func(constructor any) {
		b.RegisterBean(constructor)
	})
}

func (b *beanBuffer) Buffer() error {
	var err error
	isAlreadyInitialized := true

	b.bufferOnce.Do(func() {
		ctx := jen.NewFilePath(b.filePath)

		idx := 0
		for t, _ := range b.constructorMap {
			alias := fmt.Sprintf("i%d", idx)
			ctx.ImportAlias(t.PkgPath(), alias)
			ctx.Var().Id(fmt.Sprintf("%sBean", t.Name())).Qual(t.PkgPath(), t.Name()) // define Variable
			idx += 1
		}

		dependencyTrees := b.getDependencyTrees()
		for _, tree := range dependencyTrees {
			if e := b.registerBeanRecursive(ctx, tree); e != nil {
				err = e
				break
			}
		}

		initializeStatements := collection.Map(b.initializeOrder, func(order reflect.Type) jen.Code {
			return jen.Id(fmt.Sprintf("New%s", order.Name())).Call()
		})
		ctx.Func().Id("init").Params().Block(initializeStatements...)

		// release constructor resource
		isAlreadyInitialized = false
		b.constructorMap = nil

		// create directory if not
		if _, err := os.Stat(b.filePath); os.IsNotExist(err) {
			if err := os.MkdirAll(b.filePath, 0755); err != nil {
				panic(err)
			}
		}
		ctx.Save(filepath.Join(b.filePath, fmt.Sprintf("%s.go", b.beanSetName)))
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

func (b *beanBuffer) registerBeanRecursive(ctx *jen.File, tree collection.DescendNode[reflect.Type]) error {
	if gen, ok := b.beanChecker[tree.This()]; gen && ok {
		return nil
	}

	if tree.HasDescendants() {
		for _, child := range tree.GetDescendants() {
			err := b.registerBeanRecursive(ctx, child)
			if err != nil {
				return err
			}
		}
	}

	if constructor, ok := b.constructorMap[tree.This()]; ok {
		return b.writeCode(ctx, constructor)
	}

	return errors.New(fmt.Sprintf("fail to create bean: %s\n", tree.This().String()))
}

func (b *beanBuffer) writeCode(ctx *jen.File, constructor any) error {
	fnTypeOf := reflect.TypeOf(constructor)

	paramTypes := make([]reflect.Type, 0, fnTypeOf.NumIn())
	for i := 0; i < fnTypeOf.NumIn(); i++ {
		paramTypes = append(paramTypes, fnTypeOf.In(i))
	}

	beans := make([]string, fnTypeOf.NumIn(), fnTypeOf.NumIn())
	for i, paramType := range paramTypes {
		if gen, ok := b.beanChecker[paramType]; ok && gen {
			beans[i] = fmt.Sprintf("%sBean", paramType.Name())
		} else {
			beans[i] = "nil"
		}
	}

	if len(collection.Filter(beans, func(param string) bool { return param == "nil" })) > 0 {
		return errors.New("no bean")
	}
	generatedBeanType := fnTypeOf.Out(0)

	beanz := collection.Map(beans, func(valName string) jen.Code {
		return jen.Id(valName)
	})

	hasError := fnTypeOf.NumOut() == 2
	localBeanVariable := jen.Id("lbean")
	errStatement := jen.Line()
	if hasError {
		localBeanVariable = jen.List(jen.Id("lbean"), jen.Err())
		errStatement = jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
	}

	fmt.Println(fnTypeOf.String())
	ctx.Func().Id(fmt.Sprintf("New%s", generatedBeanType.Name())).Params().Qual(generatedBeanType.PkgPath(), generatedBeanType.Name()).
		Block(
			localBeanVariable.Op(":=").Qual(generatedBeanType.PkgPath(), fmt.Sprintf("New%s", generatedBeanType.Name())).Call(beanz...),
			errStatement,
			jen.Id(fmt.Sprintf("%sBean", generatedBeanType.Name())).Op("=").Id("lbean"),
			jen.Return(jen.Nil()),
		)

	b.beanChecker[generatedBeanType] = true
	b.initializeOrder = append(b.initializeOrder, generatedBeanType)
	return nil
}
