package intf2

import (
	"errors"
	"github.com/aivyss/bean/test/testEntity/intf1"
)

type TestInterface2 interface {
	Value2() string
}

type TestStruct2 struct{}

func (t *TestStruct2) Value2() string {
	return "test structure2"
}

type TestBeanBufferInterface2 interface{ TestBeanBufferValue2() string }

type TestBeanBufferStruct2 struct {
	interface1 intf1.TestBeanBufferInterface1
}

func (s *TestBeanBufferStruct2) TestBeanBufferValue2() string {
	return "TestBeanBufferStruct2"
}

func NewTestBeanBufferInterface2(interface1 intf1.TestBeanBufferInterface1) TestBeanBufferInterface2 {
	return &TestBeanBufferStruct2{
		interface1: interface1,
	}
}

func NewTestBeanBufferInterface2MustErr(_ intf1.TestBeanBufferInterface1) (TestBeanBufferInterface2, error) {
	return nil, errors.New("must error")
}

func NewTestBeanBufferInterface2WithNoErr(interface1 intf1.TestBeanBufferInterface1) (TestBeanBufferInterface2, error) {
	return &TestBeanBufferStruct2{
		interface1: interface1,
	}, nil
}
