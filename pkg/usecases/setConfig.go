package usecases

import (
	"fmt"
	"strings"
	"time"
)

// SetConfigInteractor wraps SetConfig operations
type SetConfigInteractor interface {
	SetConfig(userProductID int,
		config CpConfig, expiredAt time.Time) error
}

// setConfigInteractor defines the interactor for setConfig usecase
type setConfigInteractor struct {
	productRepo ProductRepository
	cacheRepo   CacheRepository
	logger      SetConfigLogger
	cacheTTL    time.Duration
}

// SetConfigLogger logs SetConfig events
type SetConfigLogger interface {
	LogErrorSettingConfig(userProductID int, err error)
	LogWarnSettingCache(userID string, err error)
}

// MakeSetConfigInteractor creates a new instance of SetConfigInteractor
func MakeSetConfigInteractor(productRepo ProductRepository,
	cacheRepo CacheRepository, logger SetConfigLogger,
	cacheTTL time.Duration) SetConfigInteractor {
	return &setConfigInteractor{productRepo: productRepo, cacheRepo: cacheRepo,
		logger: logger, cacheTTL: cacheTTL}
}

// SetConfig adds user product to repository, also sets cache
func (interactor *setConfigInteractor) SetConfig(userProductID int,
	config CpConfig, expiredAt time.Time) error {
	err := interactor.productRepo.SetExpiration(userProductID, expiredAt)
	if err != nil {
		interactor.logger.LogErrorSettingConfig(userProductID, err)
		return fmt.Errorf("cannot set control-panel partial configuration: %+v", err)
	}
	err = interactor.productRepo.SetConfig(userProductID, config)
	if err != nil {
		interactor.logger.LogErrorSettingConfig(userProductID, err)
		return fmt.Errorf("cannot set control-panel partial configuration: %+v", err)
	}
	product, err := interactor.productRepo.GetUserProductByID(userProductID)
	if err != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, err)
	}
	interactor.refreshCache(product)
	return nil
}

func (interactor *setConfigInteractor) refreshCache(product Product) {
	cacheError := interactor.cacheRepo.
		SetCache(strings.Join([]string{"user", product.UserID, string(product.Type)}, ":"),
			ProductCacheType, product, interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, cacheError)
	}
}
