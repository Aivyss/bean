package intf3

import (
	"github.com/aivyss/bean/test/testEntity/intf1"
	"github.com/aivyss/bean/test/testEntity/intf2"
)

type TestInterface3 interface {
	Value3() string
}

type TestStruct3 struct {
	interface2 intf2.TestInterface2
}

func (t *TestStruct3) Value3() string {
	return "test structure3"
}

func NewTestInterface3(interface2 intf2.TestInterface2) TestInterface3 {
	return &TestStruct3{
		interface2: interface2,
	}
}

type TestBeanBufferInterface3 interface{ TestBeanBufferValue3() string }

type TestBeanBufferStruct3 struct {
	interface1 intf1.TestBeanBufferInterface1
	interface2 intf2.TestBeanBufferInterface2
}

func (s *TestBeanBufferStruct3) TestBeanBufferValue3() string {
	return "TestBeanBufferStruct3"
}

func NewTestBeanBufferInterface3(
	interface1 intf1.TestBeanBufferInterface1,
	interface2 intf2.TestBeanBufferInterface2,
) TestBeanBufferInterface3 {
	return &TestBeanBufferStruct3{
		interface1: interface1,
		interface2: interface2,
	}
}
