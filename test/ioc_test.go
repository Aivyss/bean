package test

import (
	"errors"
	"github.com/aivyss/bean"
	"github.com/aivyss/bean/test/testEntity"
	rec "github.com/aivyss/typex/recover"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustGetBean(t *testing.T) {
	t.Run("panic occurs", func(t *testing.T) {
		err := rec.CatchPanic(func() error {
			_ = bean.MustGetBean[testEntity.TestInterface2]()

			return nil
		})
		assert.NotNil(t, err)

		bean.Clean()
	})
}

func TestRegisterBean(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		err := bean.RegisterBean(func() testEntity.TestInterface1 {
			return &testEntity.TestStruct1{}
		})
		assert.Nil(t, err)

		b, err := bean.GetBean[testEntity.TestInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, b)
		bean.Clean()
	})

	t.Run("basic test nil interface", func(t *testing.T) {
		err := bean.RegisterBean(func() testEntity.TestInterface1 {
			var i testEntity.TestInterface1 // nil
			return i
		})
		assert.Nil(t, err)
		b, err := bean.GetBean[testEntity.TestInterface1]()
		assert.Nil(t, err)
		assert.Nil(t, b)
		bean.Clean()
	})

	t.Run("nested test", func(t *testing.T) {
		err := bean.RegisterBean(func() testEntity.TestInterface2 {
			return &testEntity.TestStruct2{}
		})
		assert.Nil(t, err)

		err = bean.RegisterBean(testEntity.NewTestInterface3)
		assert.Nil(t, err)

		bean2, err := bean.GetBean[testEntity.TestInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[testEntity.TestInterface3]()
		assert.Nil(t, err)
		assert.NotNil(t, bean3)
		bean.Clean()
	})

	t.Run("error check", func(t *testing.T) {
		errMsg := "test error"
		err := bean.RegisterBean(func() (testEntity.TestInterface2, error) {
			return nil, errors.New(errMsg)
		})

		assert.Error(t, err)
		assert.Equal(t, err.Error(), errMsg)

		bean.Clean()

		err = bean.RegisterBean(func() (testEntity.TestInterface2, error) {
			return &testEntity.TestStruct2{}, nil
		})
		assert.Nil(t, err)
		b, err := bean.GetBean[testEntity.TestInterface2]()
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
		buff.RegisterBean(testEntity.NewTestBeanBufferInterface3)
		buff.RegisterBean(testEntity.NewTestBeanBufferInterface2)
		buff.RegisterBean(testEntity.NewTestBeanBufferInterface1)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		bean1, err := bean.GetBean[testEntity.TestBeanBufferInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, bean1)

		bean2, err := bean.GetBean[testEntity.TestBeanBufferInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[testEntity.TestBeanBufferInterface3]()
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
			testEntity.NewTestBeanBufferInterface3,
			testEntity.NewTestBeanBufferInterface2,
			testEntity.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		bean1, err := bean.GetBean[testEntity.TestBeanBufferInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, bean1)

		bean2, err := bean.GetBean[testEntity.TestBeanBufferInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[testEntity.TestBeanBufferInterface3]()
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
			testEntity.NewTestBeanBufferInterface3,
			testEntity.NewTestBeanBufferInterface2WithNoErr,
			testEntity.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		bean1, err := bean.GetBean[testEntity.TestBeanBufferInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, bean1)

		bean2, err := bean.GetBean[testEntity.TestBeanBufferInterface2]()
		assert.Nil(t, err)
		assert.NotNil(t, bean2)

		bean3, err := bean.GetBean[testEntity.TestBeanBufferInterface3]()
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
			testEntity.NewTestBeanBufferInterface3,
			testEntity.NewTestBeanBufferInterface2MustErr,
			testEntity.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.NotNil(t, errs)
		bean.Clean()
	})
}
