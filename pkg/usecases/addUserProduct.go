package usecases

import (
	"fmt"
	"strings"
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// AddUserProductInteractor wraps AddUserProduct operations
type AddUserProductInteractor interface {
	AddUserProduct(userID, email string,
		purchaseNumber, purchasePrice int, purchaseType domain.PurchaseType,
		productType domain.ProductType, expiredAt time.Time,
		config domain.ProductParams) error
}

// addUserProductInteractor defines the interactor for addUserProduct usecase
type addUserProductInteractor struct {
	productRepo  ProductRepository
	purchaseRepo PurchaseRepository
	cacheRepo    CacheRepository
	logger       AddUserProductLogger
	cacheTTL     time.Duration
}

// AddUserProductLogger logs AddUserProduct events
type AddUserProductLogger interface {
	LogErrorAddingProduct(userID string, err error)
	LogWarnSettingCache(userID string, err error)
}

// MakeAddUserProductInteractor creates a new instance of AddUserProductInteractor
func MakeAddUserProductInteractor(productRepo ProductRepository,
	purchaseRepo PurchaseRepository,
	cacheRepo CacheRepository, logger AddUserProductLogger,
	cacheTTL time.Duration) AddUserProductInteractor {
	return &addUserProductInteractor{productRepo: productRepo,
		purchaseRepo: purchaseRepo, cacheRepo: cacheRepo,
		logger: logger, cacheTTL: cacheTTL}
}

// AddUserProduct associates a new product to user
func (interactor *addUserProductInteractor) AddUserProduct(userID, email string,
	purchaseNumber, purchasePrice int, purchaseType domain.PurchaseType,
	productType domain.ProductType, expiredAt time.Time,
	config domain.ProductParams) error {
	err := interactor.validate(userID, productType)
	if err != nil {
		interactor.logger.LogErrorAddingProduct(userID, err)
		return err
	}
	purchase, err := interactor.purchaseRepo.CreatePurchase(purchaseNumber,
		purchasePrice, purchaseType)
	if err != nil {
		interactor.logger.LogErrorAddingProduct(userID, err)
		return fmt.Errorf("cannot create purchase: %+v", err)
	}
	product, err := interactor.productRepo.CreateUserProduct(userID, email,
		purchase, productType, expiredAt, config)
	if err != nil {
		interactor.logger.LogErrorAddingProduct(userID, err)
		return fmt.Errorf("cannot set control-panel configuration: %+v", err)
	}
	product.Purchase, err = interactor.purchaseRepo.AcceptPurchase(product.Purchase)
	if err != nil {
		interactor.logger.LogErrorAddingProduct(userID, err)
		return fmt.Errorf("cannot set control-panel configuration: %+v", err)
	}
	interactor.refreshCache(product)
	return nil
}

// validate validates conditions to create a product for user
func (interactor *addUserProductInteractor) validate(userID string,
	productType domain.ProductType) error {
	_, err := interactor.productRepo.GetUserActiveProduct(userID, productType)
	if err != ErrProductNotFound {
		if err == nil {
			err = fmt.Errorf("user already has an active product %s",
				productType)
		}
		return err
	}
	return nil
}

// refreshCache updates cache in repository for user product
func (interactor *addUserProductInteractor) refreshCache(product domain.Product) {
	cacheError := interactor.cacheRepo.
		SetCache(strings.Join([]string{"user",
			product.UserID, string(domain.PremiumCarousel)}, ":"),
			ProductCacheType, product, interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, cacheError)
	}
}
