package usecases

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

// AddUserProductInteractor wraps AddUserProduct operations
type AddUserProductInteractor interface {
	AddUserProduct(userID int, email string,
		purchaseNumber, purchasePrice int, purchaseType domain.PurchaseType,
		productType domain.ProductType, expiredAt time.Time,
		config domain.ProductParams) error
}

// addUserProductInteractor defines the interactor for addUserProduct usecase
type addUserProductInteractor struct {
	productRepo          ProductRepository
	purchaseRepo         PurchaseRepository
	cacheRepo            CacheRepository
	logger               AddUserProductLogger
	cacheTTL             time.Duration
	backendEventsRepo    BackendEventsRepository
	backendEventsEnabled bool
}

// AddUserProductLogger logs AddUserProduct events
type AddUserProductLogger interface {
	LogErrorAddingProduct(userID int, err error)
	LogWarnSettingCache(userID int, err error)
	LogWarnPushingEvent(productID int, err error)
}

// MakeAddUserProductInteractor creates a new instance of AddUserProductInteractor
func MakeAddUserProductInteractor(productRepo ProductRepository,
	purchaseRepo PurchaseRepository,
	cacheRepo CacheRepository, logger AddUserProductLogger,
	cacheTTL time.Duration, BackendEventsRepo BackendEventsRepository,
	backendEventsEnabled bool) AddUserProductInteractor {
	return &addUserProductInteractor{productRepo: productRepo,
		purchaseRepo: purchaseRepo, cacheRepo: cacheRepo,
		logger: logger, cacheTTL: cacheTTL,
		backendEventsRepo:    BackendEventsRepo,
		backendEventsEnabled: backendEventsEnabled}
}

// AddUserProduct associates a new product to user
func (interactor *addUserProductInteractor) AddUserProduct(userID int, email string,
	purchaseNumber, purchasePrice int, purchaseType domain.PurchaseType,
	productType domain.ProductType, expiredAt time.Time,
	config domain.ProductParams) error {
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
	if interactor.backendEventsEnabled {
		if err := interactor.backendEventsRepo.
			PushSoldProduct(product); err != nil {
			interactor.logger.LogWarnPushingEvent(product.ID, err)
		}
	}
	return nil
}

// refreshCache updates cache in repository for user product
func (interactor *addUserProductInteractor) refreshCache(product domain.Product) {
	cacheError := interactor.cacheRepo.
		SetCache(strings.Join([]string{"user",
			strconv.Itoa(product.UserID), string(domain.PremiumCarousel)}, ":"),
			ProductCacheType, product, interactor.cacheTTL)
	if cacheError != nil {
		interactor.logger.LogWarnSettingCache(product.UserID, cacheError)
	}
}
