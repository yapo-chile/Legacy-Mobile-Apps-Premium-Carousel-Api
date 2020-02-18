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
	LogErrorGettingUserProducts(email string, err error)
}

// MakeGetUserProductsInteractor creates a new instance of GetUserProductsInteractor
func MakeGetUserProductsInteractor(productRepo ProductRepository,
	logger GetUserProductsLogger) GetUserProductsInteractor {
	return &getUserProductsInteractor{productRepo: productRepo, logger: logger}
}

// GetUserProducts gets all user products using pagination
func (interactor *getUserProductsInteractor) GetUserProducts(email string,
	page int) ([]Product, int, int, error) {
	products, currentPage, totalPages, err := interactor.productRepo.
		GetUserProducts(email, page)
	if err != nil {
		interactor.logger.LogErrorGettingUserProducts(email, err)
		return []Product{}, 0, 0, fmt.Errorf("error loading products: %+v", err)
	}
	return products, currentPage, totalPages, nil
}
