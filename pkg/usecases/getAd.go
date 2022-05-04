package usecases

import (
	"encoding/json"
	"strings"
	"time"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
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
	ad, cacheError := interactor.getCache(listID)
	if cacheError == nil {
		return ad, nil
	}
	interactor.logger.LogWarnGettingCache(listID, cacheError)
	ad, err = interactor.adRepo.GetAd(listID)
	if err != nil {
		interactor.logger.LogErrorGettingAd(listID, err)
		return domain.Ad{}, err
	}
	interactor.refreshCache(ad)
	return ad, nil
}

func (interactor *getAdInteractor) getCache(listID string) (ad domain.Ad, cacheError error) {
	rawCachedAd, cacheError := interactor.cacheRepo.GetCache(
		strings.Join([]string{"ad", listID}, ":"), MinifiedAdDataType)
	if cacheError == nil {
		cacheError = json.Unmarshal(rawCachedAd, &ad)
	}
	if cacheError != nil {
		return domain.Ad{}, cacheError
	}
	return ad, nil
}

func (interactor *getAdInteractor) refreshCache(ad domain.Ad) {
	cacheError := interactor.cacheRepo.SetCache(
		strings.Join([]string{"ad", ad.ID}, ":"),
		MinifiedAdDataType, ad, interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(ad.ID, cacheError)
	}
}
