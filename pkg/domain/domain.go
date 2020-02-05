package domain

// Ad struct that defines ad data
type Ad struct {
	// ID defines the ListID
	ID string
	// CategoryID defines the CategoryID
	CategoryID string
	// Subject defines the ad title
	Subject string
	// Price represents the ad price
	Price float64
	// UnitOfAccount represents the ad price in a nominal monetary measure defined by
	// law. Example: for Chile is UF.
	// If unitOFAccount is defined, then it will replace the original price.
	UnitOfAccount float64
	// Currency is the symbol displayed by widget, If unitOfAccount is defined
	// then currency must be the unitOfAccount symbol
	Currency string
	URL      string
	// Image defines ads images displayed by widget
	Image    Image
	Metadata Metadata
	IsRandom bool
}

// Image struct that defines the internal structure of ad images
type Image struct {
	Full   string
	Medium string
	Small  string
}

// Metadata contains results metadata
type Metadata struct {
	Filtered int
	Total    int
}

// Ads struct that defines a group of ads
type Ads []Ad
