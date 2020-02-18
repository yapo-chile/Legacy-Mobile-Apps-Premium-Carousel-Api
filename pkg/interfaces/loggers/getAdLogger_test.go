package loggers

import (
	"testing"
)

func TestGetAdLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetAdLogger(m)
	l.LogWarnGettingCache("", nil)
	l.LogWarnSettingCache("", nil)
	l.LogErrorGettingAd("", nil)
	m.AssertExpectations(t)
}
