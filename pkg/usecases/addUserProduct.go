package usecases

import (
	"fmt"
	"strings"
	"time"
)

// GetUserAdsInteractor allows GetUserAds operations
type AddUserProductInteractor interface {
	AddUserProduct(userID, email, comment string,
		productType ProductType, expiredAt time.Time, config CpConfig) error
}

// getUserAdsInteractor defines the interactor for GetUserAds usecase
type addUserProductInteractor struct {
	productRepo ProductRepository
	cacheRepo   CacheRepository
	logger      AddUserProductLogger
}

// AddUserProductLogger logs AddUserProduct events
type AddUserProductLogger interface {
	LogErrorAddingProduct(userID string, err error)
	LogWarnSettingCache(userID string, err error)
}

// MakeAddUserProductInteractor creates a new instance of AddUserProductInteractor
func MakeAddUserProductInteractor(productRepo ProductRepository,
	cacheRepo CacheRepository, logger AddUserProductLogger) AddUserProductInteractor {
	return &addUserProductInteractor{productRepo: productRepo, cacheRepo: cacheRepo,
		logger: logger}
}

// GetUserAds retrieves user ads based on configuration repository
func (interactor *addUserProductInteractor) AddUserProduct(userID, email, comment string,
	productType ProductType, expiredAt time.Time, config CpConfig) error {
	product, err := interactor.productRepo.AddUserProduct(userID, email, comment, productType,
		expiredAt, config)
	if err != nil {
		interactor.logger.LogErrorAddingProduct(userID, err)
		return fmt.Errorf("cannot set control-panel configuration: %+v", err)
	}
	cacheError := interactor.cacheRepo.
		SetCache(strings.Join([]string{"user", userID, string(PremiumCarousel)}, ":"),
			ProductCacheType, product, time.Minute*10)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(userID, cacheError)
	}
	return nil
}
