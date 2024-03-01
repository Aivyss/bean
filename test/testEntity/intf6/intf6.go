package intf6

type TestInterface6 interface {
	Value6() string
}

type TestStruct6 struct{}

func (t *TestStruct6) Value6() string {
	return "test structure6"
}
