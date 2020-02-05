package repository

import (
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

type CpConfig struct {
	handler DbHandler
}

func MakeConfigRepository(handler DbHandler) usecases.ConfigRepository {
	return &CpConfig{
		handler: handler,
	}
}

func (repo *CpConfig) GetConfig(userID string) (usecases.CpConfig, error) {
	return usecases.CpConfig{
		Sorting:            "random",
		Categories:         []string{"2020", "1000"},
		Limit:              7,
		CustomQuery:        ``,
		Exclude:            []string{"4951846"},
		PriceRangeFrom:     0,
		PriceRangeTo:       0,
		FillGapsWithRandom: false,
	}, nil
}
