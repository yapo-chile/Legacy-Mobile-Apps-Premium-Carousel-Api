package usecases

import (
	"fmt"
)

// ExpireProductsInteractor wraps ExpireProducts operations
type ExpireProductsInteractor interface {
	ExpireProducts() error
}

// expireProductsInteractor defines the interactor for ExpireProducts usecase
type expireProductsInteractor struct {
	productRepo ProductRepository
	logger      ExpireProductsLogger
}

// ExpireProductsLogger logs ExpireProducts events
type ExpireProductsLogger interface {
	LogExpireProductsError(err error)
}

// MakeExpireProductsInteractor creates a new instance of ExpireProductsInteractor
func MakeExpireProductsInteractor(
	productRepo ProductRepository,
	logger ExpireProductsLogger,
) ExpireProductsInteractor {
	return &expireProductsInteractor{
		productRepo: productRepo,
		logger:      logger,
	}
}

// ExpireProducts set expired status for all expired products
func (interactor *expireProductsInteractor) ExpireProducts() error {
	err := interactor.productRepo.ExpireProducts()
	if err != nil {
		interactor.logger.LogExpireProductsError(err)
		return fmt.Errorf("error expiring products: %+v", err)
	}
	return nil
}
