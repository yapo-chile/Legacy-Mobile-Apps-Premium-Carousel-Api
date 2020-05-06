package loggers

import (
	"testing"
)

func TestAddUserProductLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeAddUserProductLogger(m)
	l.LogErrorAddingProduct(0, nil)
	l.LogWarnSettingCache(0, nil)
	l.LogWarnPushingEvent(0, nil)
	m.AssertExpectations(t)
}
