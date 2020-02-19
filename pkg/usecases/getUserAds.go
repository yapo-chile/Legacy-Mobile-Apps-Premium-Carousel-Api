package usecases

import (
	"fmt"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// GetUserAdsInteractor allows GetUserAds operations
type GetUserAdsInteractor interface {
	GetUserAds(userID string, exclude ...string) (domain.Ads, error)
}

// ConfigRepository interface to allows config repository operations
type ConfigRepository interface {
	GetConfig(userID string) (CpConfig, error)
}

// CpConfig holds configurations to get user ads
type CpConfig struct {
	Categories         []int
	Exclude            []string
	CustomQuery        string
	Limit              int
	PriceRangeFrom     int
	PriceRangeTo       int
	FillGapsWithRandom bool
}

// getUserAdsInteractor defines the interactor for GetUserAds usecase
type getUserAdsInteractor struct {
	adRepo     AdRepository
	configRepo ConfigRepository
}

// MakeGetUserAdsInteractor creates a new instance of GetUserAdsInteractor
func MakeGetUserAdsInteractor(adRepo AdRepository, configRepo ConfigRepository) GetUserAdsInteractor {
	return &getUserAdsInteractor{adRepo: adRepo, configRepo: configRepo}
}

// GetUserAds retrieves user ads based on configuration repository
func (interactor *getUserAdsInteractor) GetUserAds(userID string, excludeListID ...string) (domain.Ads, error) {
	// TODO implement cache logic for config
	cpConfig, err := interactor.configRepo.GetConfig(userID)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve control-panel configuration: %+v", err)
	}
	cpConfig.Exclude = append(cpConfig.Exclude, excludeListID...)
	response, err := interactor.adRepo.GetUserAds(userID, cpConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve the user's ads: %+v", err)
	}
	return response, nil
}
