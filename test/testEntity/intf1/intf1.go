package intf1

type TestInterface1 interface {
	Value1() string
}

type TestStruct1 struct{}

func (t *TestStruct1) Value1() string {
	return "test structure1"
}

type TestBeanBufferInterface1 interface{ TestBeanBufferValue1() string }

type TestBeanBufferStruct1 struct{}

func (s *TestBeanBufferStruct1) TestBeanBufferValue1() string {
	return "TestBeanBufferStruct1"
}

func NewTestBeanBufferInterface1() TestBeanBufferInterface1 {
	return &TestBeanBufferStruct1{}
}
