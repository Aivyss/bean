package intf5

import "github.com/aivyss/bean/test/testEntity/intf4"

type TestInterface5 interface {
	Value5() string
}

type TestStruct5 struct {
	interface4 intf4.TestInterface4
	interface6 TestInterface6
}

func (t *TestStruct5) Value5() string {
	return "test structure3"
}
