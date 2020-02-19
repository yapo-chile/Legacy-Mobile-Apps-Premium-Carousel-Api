package repository

import (
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

// cpConfig holds connections to get CpConfig
type cpConfig struct {
	handler DbHandler
}

// MakeConfigRepository creates a new instance of ConfigRepository
func MakeConfigRepository(handler DbHandler) usecases.ConfigRepository {
	return &cpConfig{
		handler: handler,
	}
}

// GetConfig gets controlpanel configuration for an specific userID
func (repo *cpConfig) GetConfig(userID string) (usecases.CpConfig, error) {
	// TODO: load from repo
	return usecases.CpConfig{
		Categories:         []int{},
		Limit:              10,
		CustomQuery:        ``,
		Exclude:            []string{},
		PriceRangeFrom:     0,
		PriceRangeTo:       0,
		FillGapsWithRandom: true,
	}, nil
}
