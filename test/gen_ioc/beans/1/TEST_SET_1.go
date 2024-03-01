package pkg

import (
	i2 "github.com/aivyss/bean/test/testEntity/intf1"
	i1 "github.com/aivyss/bean/test/testEntity/intf2"
	i0 "github.com/aivyss/bean/test/testEntity/intf3"
)

var TestBeanBufferInterface3Bean i0.TestBeanBufferInterface3
var TestBeanBufferInterface2Bean i1.TestBeanBufferInterface2
var TestBeanBufferInterface1Bean i2.TestBeanBufferInterface1

func NewTestBeanBufferInterface1() i2.TestBeanBufferInterface1 {
	lbean := i2.NewTestBeanBufferInterface1()

	TestBeanBufferInterface1Bean = lbean
	return nil
}
func NewTestBeanBufferInterface2() i1.TestBeanBufferInterface2 {
	lbean := i1.NewTestBeanBufferInterface2(TestBeanBufferInterface1Bean)

	TestBeanBufferInterface2Bean = lbean
	return nil
}
func NewTestBeanBufferInterface3() i0.TestBeanBufferInterface3 {
	lbean := i0.NewTestBeanBufferInterface3(TestBeanBufferInterface1Bean, TestBeanBufferInterface2Bean)

	TestBeanBufferInterface3Bean = lbean
	return nil
}
func init() {
	NewTestBeanBufferInterface1()
	NewTestBeanBufferInterface2()
	NewTestBeanBufferInterface3()
}
