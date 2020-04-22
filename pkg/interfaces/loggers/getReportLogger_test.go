package loggers

import (
	"testing"
)

func TestGetReportLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetReportLogger(m)
	l.LogErrorGettingReport(nil)
	m.AssertExpectations(t)
}
