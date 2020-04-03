package domain

import "time"

// Ad struct that defines ad data
type Ad struct {
	// ID defines the ListID
	ID string
	// UserID is the seller userID
	UserID int
	// CategoryID defines the CategoryID
	CategoryID int
	// Subject defines the ad title
	Subject string
	// Price represents the ad price
	Price float64
	// Currency is the symbol displayed by widget, If unitOfAccount is defined
	// then currency must be the unitOfAccount symbol
	Currency string
	// URL defines the ad URL
	URL string
	// Image defines ads images displayed by widget
	Image Image
	// IsRelated determines if the ad is related (true) with current adview
	// or is random (false)
	IsRelated bool
}

// Image struct that defines the internal structure of ad images
type Image struct {
	Full   string
	Medium string
	Small  string
}

// Ads struct that defines a group of ads
type Ads []Ad

// Purchase holds purchase data
type Purchase struct {
	ID        int
	Number    int
	Price     int
	Type      PurchaseType
	Status    PurchaseStatus
	CreatedAt time.Time
}

// PurchaseType defines the purchase type
type PurchaseType string

const (
	// AdminPurchase defines a purchase set by admin
	AdminPurchase PurchaseType = "ADMIN"
)

// PurchaseStatus defines the purchase status
type PurchaseStatus string

const (
	// PendingPurchase defines a pending purchase
	PendingPurchase PurchaseStatus = "PENDING"
	// AcceptedPurchase defines a accepted purchase
	AcceptedPurchase PurchaseStatus = "ACCEPTED"
	// RejectedPurchase defines a rejected purchase
	RejectedPurchase PurchaseStatus = "REJECTED"
)

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
	UserID    int
	Email     string
	Purchase  Purchase
	Status    ProductStatus
	ExpiredAt time.Time
	CreatedAt time.Time
	Config    ProductParams
}

// ProductParams holds configurations to get user ads and fill carousel
type ProductParams struct {
	Categories         []int
	Exclude            []string
	Keywords           []string
	Limit              int
	PriceRange         int
	PriceFrom          int
	PriceTo            int
	FillGapsWithRandom bool
	Comment            string
}
