package domain

// Ad struct that defines ad data
type Ad struct {
	// ID defines the ListID
	ID string
	// UserID is the seller userID
	UserID string
	// CategoryID defines the CategoryID
	CategoryID string
	// Subject defines the ad title
	Subject string
	// Price represents the ad price
	Price float64
	// Currency is the symbol displayed by widget, If unitOfAccount is defined
	// then currency must be the unitOfAccount symbol
	Currency string
	URL      string
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
