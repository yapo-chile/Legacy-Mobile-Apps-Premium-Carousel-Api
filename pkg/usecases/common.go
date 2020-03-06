package usecases

import (
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

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

// ProductRepository interface to allows config repository operations
type ProductRepository interface {
	GetUserProducts(page int) ([]Product, int, int, error)
	GetUserProductsByEmail(email string, page int) ([]Product, int, int, error)
	AddUserProduct(userID, email, comment string, productType ProductType,
		expiredAt time.Time, config CpConfig) (Product, error)
	GetUserActiveProduct(userID string,
		productType ProductType) (Product, error)
	GetUserProductsTotal() (total int)
	GetUserProductsTotalByEmail(email string) (total int)
	GetUserProductByID(userProductID int) (Product, error)
	SetConfig(userProductID int, config CpConfig) error
	SetPartialConfig(userProductID int, configMap map[string]interface{}) error
	SetExpiration(userProductID int, expiredAt time.Time) error
	SetStatus(userProductID int, status ProductStatus) error
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

// CpConfig holds configurations to get user ads
type CpConfig struct {
	Categories         []int
	Exclude            []string
	CustomQuery        string
	Limit              int
	PriceRange         int
	PriceFrom          int
	PriceTo            int
	FillGapsWithRandom bool
}

// ProductType define the available products
type ProductType string

const (
	// PremiumCarousel defines the premium carousel product
	PremiumCarousel ProductType = "PREMIUM_CAROUSEL"
)

// ProductStatus defines the product status
type ProductStatus string

const (
	// InactiveProduct defines the inactive product status
	InactiveProduct ProductStatus = "INACTIVE"
	// ActiveProduct defines the active product status
	ActiveProduct ProductStatus = "ACTIVE"
	// ExpiredProduct defines the expired product status
	ExpiredProduct ProductStatus = "EXPIRED"
)

// Product holds product information and configurations
type Product struct {
	ID        int
	Type      ProductType
	UserID    string
	Email     string
	Status    ProductStatus
	ExpiredAt time.Time
	CreatedAt time.Time
	Comment   string
	Config    CpConfig
}
