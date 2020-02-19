package usecases

import (
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// GetAdInteractor allows GetAd operations
type GetAdInteractor interface {
	GetAd(listID string) (domain.Ad, error)
}

// getAdInteractor defines the interactor for getAd usecase
type getAdInteractor struct {
	adRepo AdRepository
}

// MakeGetAdInteractor creates a new instance of GetAdInteractor
func MakeGetAdInteractor(adRepo AdRepository) GetAdInteractor {
	return &getAdInteractor{adRepo: adRepo}
}

// GetAd gets ad by given listID
func (interactor *getAdInteractor) GetAd(listID string) (domain.Ad, error) {
	// TODO implement cache logic
	return interactor.adRepo.GetAd(listID)
}
