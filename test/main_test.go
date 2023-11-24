package test

import (
	"github.com/aivyss/bean"
	"testing"
)

func TestMain(m *testing.M) {
	bean.Clean()
	m.Run()
}
