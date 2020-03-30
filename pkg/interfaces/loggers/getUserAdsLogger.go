package loggers

import (
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

type getUserAdsLogger struct {
	logger Logger
}

func (l *getUserAdsLogger) LogWarnGettingCache(userID string, err error) {
	l.logger.Warn("not able to get cache for user ads: userID %s - %+v", userID, err)
}

func (l *getUserAdsLogger) LogWarnSettingCache(userID string, err error) {
	l.logger.Warn("not able to set cache for user ads: userID %s - %+v", userID, err)
}

func (l *getUserAdsLogger) LogInfoActiveProductNotFound(userID string, product domain.Product) {
	l.logger.Info("active product not found for userID: %s. Current product is: %v (id: %d)",
		userID, product.Status, product.ID)
}

func (l *getUserAdsLogger) LogInfoProductExpired(userID string, product domain.Product) {
	l.logger.Info("the requested product %d (userID: %s) is expired at %+v", userID, product.ID,
		product.ExpiredAt)
}

func (l *getUserAdsLogger) LogErrorGettingUserAdsData(userID string, err error) {
	l.logger.Error("error getting user ads data: userID %s, error: %+v", userID, err)
}

func (l *getUserAdsLogger) LogNotEnoughAds(userID string) {
	l.logger.Error("user %s does not have enough active ads", userID)
}

// MakeGetUserAdsLogger sets up a GetUserAdsLogger instrumented
// via the provided logger
func MakeGetUserAdsLogger(logger Logger) usecases.GetUserAdsLogger {
	return &getUserAdsLogger{
		logger: logger,
	}
}
