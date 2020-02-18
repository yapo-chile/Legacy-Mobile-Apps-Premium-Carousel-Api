package loggers

import (
	"testing"
)

func TestAddUserProductLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeAddUserProductLogger(m)
	l.LogErrorAddingProduct("", nil)
	l.LogWarnSettingCache("", nil)
	m.AssertExpectations(t)
}
