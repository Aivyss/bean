# Bean
If you have experience with the **Spring Framework** in the **Java ecosystem**, you'll appreciate how straightforward and convenient Java's **Inversion of Control (IoC)** and **Dependency Injection (DI)** features are. However, finding a library with similar capabilities in Golang can be challenging.

For this reason, this library that can perform functions similar to IoC and DI supported by Java's Spring Framework has been implemented. The library, also named "bean," draws its inspiration from the beans in the Spring Framework.

# Import
```shell
$ go get github.com/aivyss/bean@latest
```

# Purpose of Use
The primary purpose of using **bean** is to delegate the authority for object creation. Developers using **bean** don't need to worry about object creation, creation order, or dependency relationships.

# What is a Bean?
The term "bean", also the name of the library, refers to the objects managed by this library. Currently, the library allows only one object per type. While Java-based libraries may allow multiple beans to be registered for a single type, this complexity is intentionally avoided in the current implementation.

If you wish to register multiple instances of `TypeA`, you can work around it using the following methods:
```go
type TypeC struct {
    TypeA
}
```
```go
type TypeB TypeA
```
# Bean Registration
There are individual and bulk registration methods for beans.
## Individual Registration
```go
err := bean.RegisterBean(func() A {
    // ...
})
```
```go
err := bean.RegisterBean(func() (B, error) {
    // ...
})
```
## Bulk Registration
```go
buf := bean.GetBeanBuffer()

buf.RegisterBean(func() A {
    // ...
})
buf.RegisterBean(func() (B, error) {
    // ...
})

err := buf.Buffer()
```
Even with bulk registration, the library checks dependency relationships and creates beans in the required order. Therefore, it is not necessary to specify the order in which RegisterBean is called.
# Obtaining a Bean
```go
a, err := bean.GetBean[A]()
if err != nil {
    // ...
}

b, err := bean.GetBean[B]()
if err != nil {
    // ...
}
```
# Example
## Code
```go
type A struct {
    b B
    d D
}

type B struct {
    c C
}

type C struct {}
type D struct {}

func NewA(
    b B,
    d D,
) (A, error) {
    return A{
        b: b,
        d: d,
    }, nil
}

func NewB(c C) (B, error) {
    return B{
        c: c,
    }, nil
}

func NewC() C {
    return C{}
}

func NewD() D {
    return D{}
}
```
```go
buf := bean.GetBeanBuffer()
buf.RegisterBean(NewA)
buf.RegisterBean(NewB)
buf.RegisterBean(NewC)
buf.RegisterBean(NewD)
```
## Dependencies
```
A
├── B
│   └── C
└── D
```
## Creation Order
1. C, D
2. B
3. A