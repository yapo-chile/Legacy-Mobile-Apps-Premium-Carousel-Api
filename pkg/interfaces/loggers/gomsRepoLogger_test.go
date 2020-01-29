package loggers

import (
	"errors"
	"testing"
)

func TestGomsRepoLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGomsRepoLogger(m)
	l.LogURI("")
	l.LogRequestErr(errors.New("Error"))
	l.LogHealthcheckOK("")
}
