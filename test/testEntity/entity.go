package testEntity

type TestInterface1 interface {
	Value1() string
}

type TestStruct1 struct{}

func (t *TestStruct1) Value1() string {
	return "test structure1"
}

// -----

type TestInterface2 interface {
	Value2() string
}

type TestStruct2 struct{}

func (t *TestStruct2) Value2() string {
	return "test structure2"
}

type TestInterface3 interface {
	Value3() string
}

type TestStruct3 struct {
	interface2 TestInterface2
}

func (t *TestStruct3) Value3() string {
	return "test structure3"
}

func NewTestInterface3(interface2 TestInterface2) TestInterface3 {
	return &TestStruct3{
		interface2: interface2,
	}
}

// -----

type TestInterface4 interface {
	Value4() string
}

type TestStruct4 struct{}

func (t *TestStruct4) Value4() string {
	return "test structure2"
}

type TestInterface5 interface {
	Value5() string
}

type TestStruct5 struct {
	interface4 TestInterface4
	interface6 TestInterface6
}

func (t *TestStruct5) Value5() string {
	return "test structure3"
}

type TestInterface6 interface {
	Value6() string
}

type TestStruct6 struct{}

func (t *TestStruct6) Value6() string {
	return "test structure6"
}

func NewTestInterface5(
	interface4 TestInterface4,
	interface6 TestInterface6,
) TestInterface5 {
	return &TestStruct5{
		interface4: interface4,
		interface6: interface6,
	}
}

// ------
type TestBeanBufferInterface1 interface{ TestBeanBufferValue1() string }
type TestBeanBufferInterface2 interface{ TestBeanBufferValue2() string }
type TestBeanBufferInterface3 interface{ TestBeanBufferValue3() string }

type TestBeanBufferStruct1 struct{}

func (s *TestBeanBufferStruct1) TestBeanBufferValue1() string {
	return "TestBeanBufferStruct1"
}

func NewTestBeanBufferInterface1() TestBeanBufferInterface1 {
	return &TestBeanBufferStruct1{}
}

type TestBeanBufferStruct2 struct {
	interface1 TestBeanBufferInterface1
}

func (s *TestBeanBufferStruct2) TestBeanBufferValue2() string {
	return "TestBeanBufferStruct2"
}

func NewTestBeanBufferInterface2(interface1 TestBeanBufferInterface1) TestBeanBufferInterface2 {
	return &TestBeanBufferStruct2{
		interface1: interface1,
	}
}

type TestBeanBufferStruct3 struct {
	interface1 TestBeanBufferInterface1
	interface2 TestBeanBufferInterface2
}

func (s *TestBeanBufferStruct3) TestBeanBufferValue3() string {
	return "TestBeanBufferStruct3"
}

func NewTestBeanBufferInterface3(
	interface1 TestBeanBufferInterface1,
	interface2 TestBeanBufferInterface2,
) TestBeanBufferInterface3 {
	return &TestBeanBufferStruct3{
		interface1: interface1,
		interface2: interface2,
	}
}
