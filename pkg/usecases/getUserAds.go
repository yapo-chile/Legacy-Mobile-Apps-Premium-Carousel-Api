package usecases

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// GetUserAdsInteractor wraps GetUserAds operations
type GetUserAdsInteractor interface {
	GetUserAds(userID string, exclude ...string) (domain.Ads, error)
}

// getUserAdsInteractor defines the interactor for GetUserAds usecase
type getUserAdsInteractor struct {
	adRepo      AdRepository
	productRepo ProductRepository
	cacheRepo   CacheRepository
	logger      GetUserAdsLogger
	cacheTTL    time.Duration
}

// GetUserAdsLogger logs getUserAds events
type GetUserAdsLogger interface {
	LogWarnGettingCache(userID string, err error)
	LogWarnSettingCache(userID string, err error)
	LogInfoActiveProductNotFound(userID string, product Product)
	LogInfoProductExpired(userID string, product Product)
	LogErrorGettingUserAdsData(userID string, err error)
}

// MakeGetUserAdsInteractor creates a new instance of GetUserAdsInteractor
func MakeGetUserAdsInteractor(adRepo AdRepository, productRepo ProductRepository,
	cacheRepo CacheRepository, logger GetUserAdsLogger, cacheTTL time.Duration) GetUserAdsInteractor {
	return &getUserAdsInteractor{adRepo: adRepo, productRepo: productRepo,
		cacheRepo: cacheRepo, logger: logger, cacheTTL: cacheTTL}
}

// GetUserAds retrieves user ads based on product configurations
func (interactor *getUserAdsInteractor) GetUserAds(userID string, excludeListID ...string) (domain.Ads, error) {
	var err error
	product, cacheError := interactor.getCache(userID)
	if cacheError != nil {
		product, err = interactor.productRepo.GetUserActiveProduct(userID,
			PremiumCarousel)
		if err != nil {
			product = Product{UserID: userID, Status: InactiveProduct}
		}
		interactor.refreshCache(product)
	}
	if product.Status != ActiveProduct {
		interactor.logger.LogInfoActiveProductNotFound(userID, product)
		return domain.Ads{}, fmt.Errorf("Product %v for user %s", product.Status, userID)
	}
	if product.ExpiredAt.Before(time.Now()) {
		product.Status = ExpiredProduct
		interactor.logger.LogInfoProductExpired(userID, product)
		interactor.refreshCache(product)
		return domain.Ads{}, interactor.productRepo.SetStatus(product.ID, product.Status)
	}
	product.Config.Exclude = append(product.Config.Exclude, excludeListID...)
	response, err := interactor.adRepo.GetUserAds(userID, product.Config)
	if err != nil {
		interactor.logger.LogErrorGettingUserAdsData(userID, err)
		return nil, fmt.Errorf("cannot retrieve the user's ads: %+v", err)
	}
	return response, nil
}

func (interactor *getUserAdsInteractor) refreshCache(product Product) {
	cacheError := interactor.cacheRepo.SetCache(
		strings.Join([]string{"user", product.UserID, string(PremiumCarousel)}, ":"),
		ProductCacheType,
		product,
		interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, cacheError)
	}
}

func (interactor *getUserAdsInteractor) getCache(userID string) (product Product,
	cacheError error) {
	rawCachedProduct, cacheError := interactor.cacheRepo.GetCache(
		strings.Join([]string{"user", userID, string(PremiumCarousel)}, ":"),
		ProductCacheType)
	if cacheError == nil {
		cacheError = json.Unmarshal(rawCachedProduct, &product)
	}
	if cacheError != nil {
		interactor.logger.LogWarnGettingCache(userID, cacheError)
		return Product{}, cacheError
	}
	return product, nil
}
