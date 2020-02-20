package usecases

import (
	"fmt"
	"strings"
	"time"
)

// AddUserProductInteractor wraps AddUserProduct operations
type AddUserProductInteractor interface {
	AddUserProduct(userID, email, comment string,
		productType ProductType, expiredAt time.Time, config CpConfig) error
}

// addUserProductInteractor defines the interactor for addUserProduct usecase
type addUserProductInteractor struct {
	productRepo ProductRepository
	cacheRepo   CacheRepository
	logger      AddUserProductLogger
	cacheTTL    time.Duration
}

// AddUserProductLogger logs AddUserProduct events
type AddUserProductLogger interface {
	LogErrorAddingProduct(userID string, err error)
	LogWarnSettingCache(userID string, err error)
}

// MakeAddUserProductInteractor creates a new instance of AddUserProductInteractor
func MakeAddUserProductInteractor(productRepo ProductRepository,
	cacheRepo CacheRepository, logger AddUserProductLogger,
	cacheTTL time.Duration) AddUserProductInteractor {
	return &addUserProductInteractor{productRepo: productRepo, cacheRepo: cacheRepo,
		logger: logger, cacheTTL: cacheTTL}
}

// AddUserProduct adds user product to repository, also sets cache
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
			ProductCacheType, product, interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(userID, cacheError)
	}
	return nil
}
