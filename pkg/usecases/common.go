package usecases

import (
	"errors"
	"time"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
)

// AdMedia holds ad images data
type AdMedia struct {
	// ID image unique ID
	ID int `json:"ID"`
	// SeqNo is the image sequence number to display in inblocket platform
	SeqNo int `json:"SeqNo"`
}

// Ad holds ad response from external source
type Ad struct {
	AdID          int64                `json:"adId"`
	ListID        int64                `json:"listId"`
	UserID        int64                `json:"userId"`
	Type          string               `json:"type"`
	Phone         string               `json:"phone"`
	Location      Location             `json:"location"`
	Category      Category             `json:"category"`
	Name          string               `json:"name"`
	URL           string               `json:"url"`
	Subject       string               `json:"subject"`
	Body          string               `json:"body"`
	Price         float64              `json:"price"`
	OldPrice      float64              `json:"oldPrice"`
	ListTime      time.Time            `json:"listTime"`
	Media         []AdMedia            `json:"media"`
	PublisherType string               `json:"publisherType"`
	Params        map[string]Param     `json:"params"`
}

// Param represents additional parameters on ads
type Param struct {
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Translate interface{} `json:"translate"`
}

// Category represents a Yapo category details
type Category struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	ParentID   int64  `json:"parentId"`
	ParentName string `json:"parentName"`
}

// Location represents a location object on Ad
type Location struct {
	RegionID    int64  `json:"regionId"`
	RegionName  string `json:"regionName"`
	ComunneID   int64  `json:"communeId"`
	CommuneName string `json:"communeName"`
}

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
