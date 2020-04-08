package usecases

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// SetPartialConfigInteractor wraps SetPartialConfig operations
type SetPartialConfigInteractor interface {
	SetPartialConfig(userProductID int,
		configMap map[string]interface{}) error
}

// setPartialConfigInteractor defines the interactor for setPartialConfig usecase
type setPartialConfigInteractor struct {
	productRepo ProductRepository
	cacheRepo   CacheRepository
	logger      SetPartialConfigLogger
	cacheTTL    time.Duration
}

// SetPartialConfigLogger logs SetPartialConfig events
type SetPartialConfigLogger interface {
	LogErrorSettingPartialConfig(userProductID int, err error)
	LogWarnSettingCache(userID int, err error)
}

// MakeSetPartialConfigInteractor creates a new instance of SetPartialConfigInteractor
func MakeSetPartialConfigInteractor(productRepo ProductRepository,
	cacheRepo CacheRepository, logger SetPartialConfigLogger,
	cacheTTL time.Duration) SetPartialConfigInteractor {
	return &setPartialConfigInteractor{productRepo: productRepo, cacheRepo: cacheRepo,
		logger: logger, cacheTTL: cacheTTL}
}

// SetPartialConfig sets partial configuration to userProduct also sets cache
func (interactor *setPartialConfigInteractor) SetPartialConfig(userProductID int,
	configMap map[string]interface{}) error {
	err := interactor.productRepo.SetPartialConfig(userProductID, configMap)
	if err != nil {
		interactor.logger.LogErrorSettingPartialConfig(userProductID, err)
		return fmt.Errorf("cannot set control-panel partial configuration: %+v", err)
	}
	product, err := interactor.productRepo.GetUserProductByID(userProductID)
	if err != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, err)
		return err
	}
	interactor.refreshCache(product)
	return nil
}

func (interactor *setPartialConfigInteractor) refreshCache(product domain.Product) {
	cacheError := interactor.cacheRepo.
		SetCache(strings.Join([]string{"user", strconv.Itoa(product.UserID),
			string(product.Type)}, ":"), ProductCacheType, product, interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, cacheError)
	}
}
