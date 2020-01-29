package loggers

import (
	"testing"
)

func TestGetUserDataHandlerLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetUserDataHandlerLogger(m)
	l.LogBadRequest(nil)
	l.LogErrorGettingInternalData(nil)
	m.AssertExpectations(t)
}
