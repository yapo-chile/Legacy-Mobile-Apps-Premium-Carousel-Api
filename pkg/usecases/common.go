package usecases

import (
	"errors"
	"time"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
)

// GomsRepository interface that represents all the methods available to
// interact with premium-carousel-api microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}

// AdRepository allows get ads data
type AdRepository interface {
	GetUserAds(userID int,
		productParams domain.ProductParams) (domain.Ads, error)
	GetAd(listID string) (domain.Ad, error)
}

// PurchaseRepository interface to allows purchase repository operations
type PurchaseRepository interface {
	CreatePurchase(purchaseNumber, price int,
		purchaseType domain.PurchaseType) (domain.Purchase, error)
	AcceptPurchase(purchase domain.Purchase) (domain.Purchase, error)
}

// ErrProductNotFound defines error for product not found
var ErrProductNotFound error = errors.New("Product not found")

// ProductRepository interface to allows product repository operations
type ProductRepository interface {
	GetUserProducts(page int) ([]domain.Product, int, int, error)
	GetUserProductsByEmail(email string, page int) ([]domain.Product,
		int, int, error)
	CreateUserProduct(userID int, email string,
		purchase domain.Purchase, productType domain.ProductType,
		expiredAt time.Time, config domain.ProductParams) (domain.Product, error)
	GetUserActiveProduct(userID int,
		productType domain.ProductType) (domain.Product, error)
	GetUserProductsTotal() (total int)
	GetUserProductsTotalByEmail(email string) (total int)
	GetUserProductByID(userProductID int) (domain.Product, error)
	SetConfig(userProductID int, config domain.ProductParams) error
	SetPartialConfig(userProductID int, configMap map[string]interface{}) error
	SetExpiration(userProductID int, expiredAt time.Time) error
	SetStatus(userProductID int, status domain.ProductStatus) error
	GetReport(startDate, endDate time.Time) ([]domain.Product, error)
	ExpireProducts() error
}

// CacheType defines the user cache type
type CacheType string

const (
	// ProductCacheType represents a product cache type
	ProductCacheType CacheType = "cache-product"
	// MinifiedAdDataType represents a minified ad data type
	MinifiedAdDataType CacheType = "cache-minified-ad-data"
)

// CacheRepository implements cache repository operations
type CacheRepository interface {
	SetCache(key string, typ CacheType, data interface{},
		expiration time.Duration) error
	GetCache(key string, typ CacheType) ([]byte, error)
}

// BackendEventsRepository allows push events to backend events queue
type BackendEventsRepository interface {
	PushSoldProduct(product domain.Product) error
}
