package loggers

import (
	"testing"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
)

func TestGetUserAdsLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetUserAdsLogger(m)
	l.LogWarnGettingCache(0, nil)
	l.LogWarnSettingCache(0, nil)
	l.LogInfoActiveProductNotFound(0, domain.Product{})
	l.LogInfoProductExpired(0, domain.Product{})
	l.LogErrorGettingUserAdsData(0, nil)
	l.LogNotEnoughAds(0)
	m.AssertExpectations(t)
}
