package loggers

import (
	"testing"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

func TestGetUserAdsLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetUserAdsLogger(m)
	l.LogWarnGettingCache("", nil)
	l.LogWarnSettingCache("", nil)
	l.LogInfoActiveProductNotFound("", domain.Product{})
	l.LogInfoProductExpired("", domain.Product{})
	l.LogErrorGettingUserAdsData("", nil)
	l.LogNotEnoughAds("")
	m.AssertExpectations(t)
}
