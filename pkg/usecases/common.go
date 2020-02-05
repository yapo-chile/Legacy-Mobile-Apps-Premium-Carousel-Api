package usecases

import "github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"

// GomsRepository interface that represents all the methods available to
// interact with premium-carousel-api microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}

// UserBasicData is the structure that contains the basic user data
type UserBasicData struct {
	UserID  string
	Name    string
	Phone   string
	Gender  string
	Country string
	Region  string
	Commune string
}

// UserProfileRepository defines the methods that a User Profile repository should have
type UserProfileRepository interface {
	// GetUserData gets the user data based on his email
	GetUserProfileData(pSHA1E string) (UserBasicData, error)
}

// AdRepository allows get ads data
type AdRepository interface {
	GetUserAds(userID string, cpConfig CpConfig) (domain.Ads, error)
}

// SearchResponse object to recieve search response data type
type SearchResponse map[string][]struct {
	ListID        string `json:"listId"`
	Subject       string `json:"subject"`
	RegionID      string `json:"region"`
	CategoryID    string `json:"category"`
	Price         string `json:"price"`
	UnitOfAccount string `json:"priceUf"`
	Currency      string `json:"currency"`
	Image         Image  `json:"image"`
	URL           string
}

// Image struct that defines the internal structure of the images
// that search-ms retrieves
type Image struct {
	MainImage string `json:"mainImage"`
	Thumbs    string `json:"thumbs"`
	Thumbli   string `json:"thumbli"`
}

// SearchInput object to recieve search input data type
type SearchInput map[string][]interface{}
