package gen_ioc

import (
	"github.com/aivyss/bean"
	bean2 "github.com/aivyss/bean/gen"
	"github.com/aivyss/bean/test/testEntity/intf1"
	"github.com/aivyss/bean/test/testEntity/intf2"
	"github.com/aivyss/bean/test/testEntity/intf3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBeanBuffer(t *testing.T) {
	t.Run("[1] buffer test - no error", func(t *testing.T) {
		buff := bean2.GetBeanBuffer("TEST_SET_1", "./beans/1")
		buff.RegisterBean(intf3.NewTestBeanBufferInterface3)
		buff.RegisterBean(intf2.NewTestBeanBufferInterface2)
		buff.RegisterBean(intf1.NewTestBeanBufferInterface1)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		t.Run("2 times buffer", func(t *testing.T) {
			err := buff.Buffer()
			assert.NotNil(t, err)
		})
		bean.Clean()
	})

	t.Run("[2] buffer test - no error", func(t *testing.T) {
		buff := bean2.GetBeanBuffer("TEST_SET_2", "./beans/2")
		buff.RegisterBeans(
			intf3.NewTestBeanBufferInterface3,
			intf2.NewTestBeanBufferInterface2,
			intf1.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		t.Run("2 times buffer", func(t *testing.T) {
			err := buff.Buffer()
			assert.NotNil(t, err)
		})
		bean.Clean()
	})

	// BUG reflect.Type doesn't have function name at runtime...
	t.Run("[3] buffer test - no error", func(t *testing.T) {
		buff := bean2.GetBeanBuffer("TEST_SET_3", "./beans/3")
		buff.RegisterBeans(
			intf3.NewTestBeanBufferInterface3,
			intf2.NewTestBeanBufferInterface2WithNoErr,
			intf1.NewTestBeanBufferInterface1,
		)

		errs := buff.Buffer()
		assert.Empty(t, errs)

		t.Run("2 times buffer", func(t *testing.T) {
			err := buff.Buffer()
			assert.NotNil(t, err)
		})
		bean.Clean()
	})

	// BUG reflect.Type doesn't have function name at runtime...
	t.Run("buffer test - must error", func(t *testing.T) {
		buff := bean2.GetBeanBuffer("TEST_SET_4", "./beans/4")
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
