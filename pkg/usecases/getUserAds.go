package usecases

import (
	"fmt"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type GetUserAdsInteractor interface {
	GetUserAds(userID string) (domain.Ads, error)
}

type ConfigRepository interface {
	GetConfig(userID string) (CpConfig, error)
}

type CpConfig struct {
	Sorting            string
	Categories         []string
	Exclude            []string
	CustomQuery        string
	Limit              int
	PriceRangeFrom     int
	PriceRangeTo       int
	FillGapsWithRandom bool
}

type ReportRepository interface {
	Save(ads domain.Ads) error
}

// getUserAdsInteractor defines the interactor
type getUserAdsInteractor struct {
	adRepo     AdRepository
	configRepo ConfigRepository
	reportRepo ReportRepository
}

func MakeGetUserAdsInteractor(adRepo AdRepository, configRepo ConfigRepository) GetUserAdsInteractor {
	return &getUserAdsInteractor{adRepo: adRepo, configRepo: configRepo}
}

// GetUser retrieves the basic data of a user given a mail
func (interactor *getUserAdsInteractor) GetUserAds(userID string) (domain.Ads, error) {
	cpConfig, err := interactor.configRepo.GetConfig(userID)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve control-panel configuration: %+v", err)
	}
	response, err := interactor.adRepo.GetUserAds(userID, cpConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve the user's ads: %+v", err)
	}
	return response, nil
}
