# Is there a way to register more than one bean for a single type?

No, it's not possible due to technical limitations with reflection. While it is possible to introduce such functionality, it would increase complexity, so it is not currently implemented. However, there is a workaround.

```go
type TypeA ...

func NewTypeAOne() (TypeA, error) {
    // ...
    
    return one, nil
}

func NewTypeATwo() (TypeA, error) {
    // ...
    
    return two, nil
} 
```
The above example is not possible. However, the following is achievable:

```go
type TypeA ...
type TypeB TypeA

func NewTypeAOne() (TypeA, error) {
    // ...
    
    return one, nil
}

func NewTypeBTwo() (TypeB, error) {
    // ...
    
    return TypeB(two), nil
}


// ...

func SomeFunction() {
    a, err := bean.GetBean[TypeA]()
    
    // ...
    b, err := bean.GetBean[TypeB]()
}
```
By simply defining another type with a different name, you can achieve what you want. Most issues can be solved using basic Go syntax.

However, since an alias is not recognized as a different type in Go, the following approach is not possible:
```go
type TypeA ...
type TypeB = TypeA

func NewTypeAOne() (TypeA, error) {
    // ...
    
    return one, nil
}

func NewTypeBTwo() (TypeB, error) {
    // ...
    
    return TypeB(two), nil
}
```

# Are there any limitations on constructor functions?
Yes, there are.

## Input Parameters
Any type designated as a bean can be used. However, if there are parameters that are not designated as beans, the bean package will return an error. There are no restrictions on the number of parameters.

## Returns
Only two return patterns are allowed: either `(TypeA)` or `(TypeA, error)`.