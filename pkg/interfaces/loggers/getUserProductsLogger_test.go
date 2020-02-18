package loggers

import (
	"testing"
)

func TestGetUserProductsLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetUserProductsLogger(m)
	l.LogErrorGettingUserProducts("", nil)
	m.AssertExpectations(t)
}
