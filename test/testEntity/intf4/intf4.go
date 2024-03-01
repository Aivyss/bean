package intf4

type TestInterface4 interface {
	Value4() string
}

type TestStruct4 struct{}

func (t *TestStruct4) Value4() string {
	return "test structure2"
}
