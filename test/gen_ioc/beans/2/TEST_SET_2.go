package pkg

import (
	i1 "github.com/aivyss/bean/test/testEntity/intf1"
	i0 "github.com/aivyss/bean/test/testEntity/intf2"
	i2 "github.com/aivyss/bean/test/testEntity/intf3"
)

var TestBeanBufferInterface2Bean i0.TestBeanBufferInterface2
var TestBeanBufferInterface1Bean i1.TestBeanBufferInterface1
var TestBeanBufferInterface3Bean i2.TestBeanBufferInterface3

func NewTestBeanBufferInterface1() i1.TestBeanBufferInterface1 {
	lbean := i1.NewTestBeanBufferInterface1()

	TestBeanBufferInterface1Bean = lbean
	return nil
}
func NewTestBeanBufferInterface2() i0.TestBeanBufferInterface2 {
	lbean := i0.NewTestBeanBufferInterface2(TestBeanBufferInterface1Bean)

	TestBeanBufferInterface2Bean = lbean
	return nil
}
func NewTestBeanBufferInterface3() i2.TestBeanBufferInterface3 {
	lbean := i2.NewTestBeanBufferInterface3(TestBeanBufferInterface1Bean, TestBeanBufferInterface2Bean)

	TestBeanBufferInterface3Bean = lbean
	return nil
}
func init() {
	NewTestBeanBufferInterface1()
	NewTestBeanBufferInterface2()
	NewTestBeanBufferInterface3()
}
