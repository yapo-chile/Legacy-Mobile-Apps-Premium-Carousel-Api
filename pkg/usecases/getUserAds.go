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
	GetUserAds(currentAdview domain.Ad) (domain.Ads, error)
}

// getUserAdsInteractor defines the interactor for GetUserAds usecase
type getUserAdsInteractor struct {
	adRepo          AdRepository
	productRepo     ProductRepository
	cacheRepo       CacheRepository
	logger          GetUserAdsLogger
	cacheTTL        time.Duration
	minAdsToDisplay int
}

// GetUserAdsLogger logs getUserAds events
type GetUserAdsLogger interface {
	LogNotEnoughAds(userID string)
	LogWarnGettingCache(userID string, err error)
	LogWarnSettingCache(userID string, err error)
	LogErrorGettingUserAdsData(userID string, err error)
	LogInfoProductExpired(userID string, product domain.Product)
	LogInfoActiveProductNotFound(userID string, product domain.Product)
}

// MakeGetUserAdsInteractor creates a new instance of GetUserAdsInteractor
func MakeGetUserAdsInteractor(adRepo AdRepository, productRepo ProductRepository,
	cacheRepo CacheRepository, logger GetUserAdsLogger,
	cacheTTL time.Duration, minAdsToDisplay int) GetUserAdsInteractor {
	return &getUserAdsInteractor{adRepo: adRepo,
		productRepo: productRepo, cacheRepo: cacheRepo,
		logger: logger, cacheTTL: cacheTTL, minAdsToDisplay: minAdsToDisplay}
}

// GetUserAds retrieves user ads based on product configurations
func (interactor *getUserAdsInteractor) GetUserAds(currentAdview domain.Ad) (ads domain.Ads, err error) {
	userID := currentAdview.UserID
	product, cacheError := interactor.getCache(userID)
	if cacheError != nil {
		product, err = interactor.productRepo.GetUserActiveProduct(userID,
			domain.PremiumCarousel)
		if err != nil {
			product = domain.Product{UserID: userID, Status: domain.InactiveProduct}
		}
		interactor.refreshCache(product)
	}
	if product.Status != domain.ActiveProduct {
		interactor.logger.LogInfoActiveProductNotFound(userID, product)
		return domain.Ads{}, fmt.Errorf("Product %v for user %s", product.Status, userID)
	}
	if product.ExpiredAt.Before(time.Now()) {
		product.Status = domain.ExpiredProduct
		interactor.logger.LogInfoProductExpired(userID, product)
		interactor.refreshCache(product)
		return domain.Ads{}, interactor.productRepo.SetStatus(product.ID, product.Status)
	}
	product.Config.Exclude = append(product.Config.Exclude, currentAdview.ID)
	if product.Config.PriceRange > 0 {
		product.Config.PriceFrom = int(currentAdview.Price) - product.Config.PriceRange
		product.Config.PriceTo = int(currentAdview.Price) + product.Config.PriceRange
	}
	ads, err = interactor.adRepo.GetUserAds(userID, product.Config)
	if err != nil {
		interactor.logger.LogErrorGettingUserAdsData(userID, err)
		return domain.Ads{}, fmt.Errorf("cannot retrieve the user's ads: %+v", err)
	}
	if interactor.minAdsToDisplay > 0 && len(ads) < interactor.minAdsToDisplay {
		interactor.logger.LogNotEnoughAds(userID)
		return domain.Ads{}, fmt.Errorf("user %s does not have enough active ads", userID)
	}
	return ads, nil
}

func (interactor *getUserAdsInteractor) refreshCache(product domain.Product) {
	cacheError := interactor.cacheRepo.SetCache(
		strings.Join([]string{"user", product.UserID,
			string(domain.PremiumCarousel)}, ":"),
		ProductCacheType,
		product,
		interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, cacheError)
	}
}

func (interactor *getUserAdsInteractor) getCache(userID string) (product domain.Product,
	cacheError error) {
	rawCachedProduct, cacheError := interactor.cacheRepo.GetCache(
		strings.Join([]string{"user", userID,
			string(domain.PremiumCarousel)}, ":"),
		ProductCacheType)
	if cacheError == nil {
		cacheError = json.Unmarshal(rawCachedProduct, &product)
	}
	if cacheError != nil {
		interactor.logger.LogWarnGettingCache(userID, cacheError)
		return domain.Product{}, cacheError
	}
	return product, nil
}
