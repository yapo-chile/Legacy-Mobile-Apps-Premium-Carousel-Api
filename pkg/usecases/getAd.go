package usecases

import (
	"encoding/json"
	"strings"
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// GetAdInteractor wraps GetAd operations
type GetAdInteractor interface {
	GetAd(listID string) (domain.Ad, error)
}

// getAdInteractor defines the interactor for getAd usecase
type getAdInteractor struct {
	adRepo    AdRepository
	cacheRepo CacheRepository
	logger    GetAdLogger
	cacheTTL  time.Duration
}

// GetAdLogger logs GetAd events
type GetAdLogger interface {
	LogWarnGettingCache(listID string, err error)
	LogWarnSettingCache(listID string, err error)
	LogErrorGettingAd(listID string, err error)
}

// MakeGetAdInteractor creates a new instance of GetAdInteractor
func MakeGetAdInteractor(adRepo AdRepository,
	cacheRepo CacheRepository, logger GetAdLogger,
	cacheTTL time.Duration) GetAdInteractor {
	return &getAdInteractor{adRepo: adRepo, cacheRepo: cacheRepo,
		logger: logger, cacheTTL: cacheTTL}
}

// GetAd gets ad by given listID
func (interactor *getAdInteractor) GetAd(listID string) (ad domain.Ad, err error) {
	rawCachedAd, cacheError := interactor.cacheRepo.GetCache(
		strings.Join([]string{"ad", listID}, ":"), MinifiedAdDataType)
	if cacheError == nil {
		cacheError = json.Unmarshal(rawCachedAd, &ad)
	}
	if cacheError == nil {
		return ad, nil
	}
	interactor.logger.LogWarnGettingCache(listID, cacheError)
	ad, err = interactor.adRepo.GetAd(listID)
	if err != nil {
		interactor.logger.LogErrorGettingAd(listID, err)
		return domain.Ad{}, err
	}
	cacheError = interactor.cacheRepo.SetCache(
		strings.Join([]string{"ad", listID}, ":"),
		MinifiedAdDataType, ad, interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(listID, cacheError)
	}
	return ad, nil
}
