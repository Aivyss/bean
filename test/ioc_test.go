package test

import (
	"errors"
	"github.com/aivyss/bean"
	"github.com/aivyss/bean/test/testEntity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterBean(t *testing.T) {
	t.Run("basic test", func(t *testing.T) {
		err := bean.RegisterBean(func() testEntity.TestInterface1 {
			return &testEntity.TestStruct1{}
		})

		assert.Nil(t, err)

		b, err := bean.GetBean[testEntity.TestInterface1]()
		assert.Nil(t, err)
		assert.NotNil(t, b)
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
}

func TestBeanBuffer(t *testing.T) {
	t.Run("buffer test - no error", func(t *testing.T) {
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
	})
}

func TestBuffer2(t *testing.T) {

}
