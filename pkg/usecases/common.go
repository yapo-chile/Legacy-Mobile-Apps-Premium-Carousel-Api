package usecases

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"

// GomsRepository interface that represents all the methods available to
// interact with premium-carousel-api microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}

// AdRepository allows get ads data
type AdRepository interface {
	GetUserAds(userID string, cpConfig CpConfig) (domain.Ads, error)
	GetAd(listID string) (domain.Ad, error)
}
