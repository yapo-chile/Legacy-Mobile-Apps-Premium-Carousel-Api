package loggers

import (
	"testing"
)

func TestProductRepoLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeProductRepositoryLogger(m)
	l.LogWarnPartialConfigNotSupported("", "")
	m.AssertExpectations(t)
}
