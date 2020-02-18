package loggers

import (
	"testing"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

func TestGetUserAdsLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetUserAdsLogger(m)
	l.LogWarnGettingCache("", nil)
	l.LogWarnSettingCache("", nil)
	l.LogInfoActiveProductNotFound("")
	l.LogInfoProductExpired("", usecases.Product{})
	l.LogErrorGettingUserAdsData("", nil)
	m.AssertExpectations(t)
}
