package loggers

import (
	"testing"
)

func TestSetPartialConfigLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeSetPartialConfigLogger(m)
	l.LogWarnSettingCache("", nil)
	l.LogErrorSettingPartialConfig(1, nil)
	m.AssertExpectations(t)
}
