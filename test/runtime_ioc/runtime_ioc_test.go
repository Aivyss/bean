package runtime_ioc

import (
	"errors"
	"github.com/aivyss/bean"
	"github.com/aivyss/bean/test/testEntity/intf1"
	"github.com/aivyss/bean/test/testEntity/intf2"
	"github.com/aivyss/bean/test/testEntity/intf3"
	rec "github.com/aivyss/typex/recover"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustGetBean(t *testing.T) {
	t.Run("panic occurs", func(t *testing.T) {
		err := rec.CatchPanic(func() error {
			_ = bean.MustGetBean[intf2.TestInterface2]()

			return nil
		})
		assert.NotNil(t, err)

		bean.Clean()
	})
}

func TestRegisterBean(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		err := bean.RegisterBean(func() intf1.TestInterface1 {
			return &intf1.TestStruct1{}
		})
		assert.Nil(t, err)

		b, err := bean.GetBean[intf1.TestInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, b)
		bean.Clean()
	})

	t.Run("basic test nil interface", func(t *testing.T) {
		err := bean.RegisterBean(func() intf1.TestInterface1 {
			var i intf1.TestInterface1 // nil
			return i
		})
		assert.Nil(t, err)
		b, err := bean.GetBean[intf1.TestInterface1]()
		assert.Nil(t, err)
		assert.Nil(t, b)
		bean.Clean()
	})

	t.Run("nested test", func(t *testing.T) {
		err := bean.RegisterBean(func() intf2.TestInterface2 {
			return &intf2.TestStruct2{}
		})
		assert.Nil(t, err)

		err = bean.RegisterBean(intf3.NewTestInterface3)
		assert.Nil(t, err)

		bean2, err := bean.GetBean[intf2.TestInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[intf3.TestInterface3]()
		assert.Nil(t, err)
		assert.NotNil(t, bean3)
		bean.Clean()
	})

	t.Run("error check", func(t *testing.T) {
		errMsg := "test error"
		err := bean.RegisterBean(func() (intf2.TestInterface2, error) {
			return nil, errors.New(errMsg)
		})

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errMsg)

		bean.Clean()

		err = bean.RegisterBean(func() (intf2.TestInterface2, error) {
			return &intf2.TestStruct2{}, nil
		})
		assert.Nil(t, err)
		b, err := bean.GetBean[intf2.TestInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, b)
	})
}

func TestBasicTestPointer(t *testing.T) {
	type TestStruct struct{}
	newTestStruct := func() *TestStruct { return &TestStruct{} }
	err := bean.RegisterBean(newTestStruct)
	assert.Nil(t, err)

	b, err := bean.GetBean[*TestStruct]()
	assert.Nil(t, err)
	assert.NotNil(t, b)
	bean.Clean()
}

func TestBeanBuffer(t *testing.T) {
	t.Run("[1] buffer test - no error", func(t *testing.T) {
		buff := bean.GetBeanBuffer()
		buff.RegisterBean(intf3.NewTestBeanBufferInterface3)
		buff.RegisterBean(intf2.NewTestBeanBufferInterface2)
		buff.RegisterBean(intf1.NewTestBeanBufferInterface1)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		bean1, err := bean.GetBean[intf1.TestBeanBufferInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, bean1)

		bean2, err := bean.GetBean[intf2.TestBeanBufferInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[intf3.TestBeanBufferInterface3]()
		assert.Nil(t, err)
		assert.NotNil(t, bean3)

		t.Run("2 times buffer", func(t *testing.T) {
			err := buff.Buffer()
			assert.NotNil(t, err)
		})
		bean.Clean()
	})

	t.Run("[2] buffer test - no error", func(t *testing.T) {
		buff := bean.GetBeanBuffer()
		buff.RegisterBeans(
			intf3.NewTestBeanBufferInterface3,
			intf2.NewTestBeanBufferInterface2,
			intf1.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		bean1, err := bean.GetBean[intf1.TestBeanBufferInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, bean1)

		bean2, err := bean.GetBean[intf2.TestBeanBufferInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[intf3.TestBeanBufferInterface3]()
		assert.Nil(t, err)
		assert.NotNil(t, bean3)

		t.Run("2 times buffer", func(t *testing.T) {
			err := buff.Buffer()
			assert.NotNil(t, err)
		})
		bean.Clean()
	})

	t.Run("[3] buffer test - no error", func(t *testing.T) {
		buff := bean.GetBeanBuffer()
		buff.RegisterBeans(
			intf3.NewTestBeanBufferInterface3,
			intf2.NewTestBeanBufferInterface2WithNoErr,
			intf1.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		bean1, err := bean.GetBean[intf1.TestBeanBufferInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, bean1)

		bean2, err := bean.GetBean[intf2.TestBeanBufferInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[intf3.TestBeanBufferInterface3]()
		assert.Nil(t, err)
		assert.NotNil(t, bean3)

		t.Run("2 times buffer", func(t *testing.T) {
			err := buff.Buffer()
			assert.NotNil(t, err)
		})
		bean.Clean()
	})

	t.Run("buffer test - must error", func(t *testing.T) {
		buff := bean.GetBeanBuffer()
		buff.RegisterBeans(
			intf3.NewTestBeanBufferInterface3,
			intf2.NewTestBeanBufferInterface2MustErr,
			intf1.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.NotNil(t, errs)
		bean.Clean()
	})
}
