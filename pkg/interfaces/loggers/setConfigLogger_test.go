package loggers

import (
	"testing"
)

func TestSetConfigLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeSetConfigLogger(m)
	l.LogErrorSettingConfig(1, nil)
	l.LogWarnSettingCache(0, nil)
	m.AssertExpectations(t)
}
