package loggers

import (
	"testing"
)

func TestExpireProductsLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeExpireProductsLogger(m)
	l.LogExpireProductsError(nil)
	m.AssertExpectations(t)
}
