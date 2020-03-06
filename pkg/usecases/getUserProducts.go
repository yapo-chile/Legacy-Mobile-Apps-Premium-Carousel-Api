package usecases

import (
	"fmt"
)

// GetUserProductsInteractor wraps GetUserProducts operations
type GetUserProductsInteractor interface {
	GetUserProducts(email string, page int) ([]Product, int, int, error)
}

// getUserProductsInteractor defines the interactor for GetUserProducts usecase
type getUserProductsInteractor struct {
	productRepo ProductRepository
	logger      GetUserProductsLogger
}

// GetUserProductsLogger logs GetUserProducts events
type GetUserProductsLogger interface {
	LogErrorGettingUserProducts(err error)
	LogErrorGettingUserProductsByEmail(email string, err error)
}

// MakeGetUserProductsInteractor creates a new instance of GetUserProductsInteractor
func MakeGetUserProductsInteractor(productRepo ProductRepository,
	logger GetUserProductsLogger) GetUserProductsInteractor {
	return &getUserProductsInteractor{productRepo: productRepo, logger: logger}
}

// GetUserProducts gets all user products using pagination
func (interactor *getUserProductsInteractor) GetUserProducts(email string,
	page int) (products []Product, currentPage int, totalPages int, err error) {
	if email == "" {
		products, currentPage, totalPages, err = interactor.productRepo.
			GetUserProducts(page)
		if err != nil {
			interactor.logger.LogErrorGettingUserProducts(err)
			return []Product{}, 0, 0, fmt.Errorf("error loading products: %+v", err)
		}
	} else {
		products, currentPage, totalPages, err = interactor.productRepo.
			GetUserProductsByEmail(email, page)
		if err != nil {
			interactor.logger.LogErrorGettingUserProductsByEmail(email, err)
			return []Product{}, 0, 0, fmt.Errorf("error loading products: %+v", err)
		}
	}
	return
}
