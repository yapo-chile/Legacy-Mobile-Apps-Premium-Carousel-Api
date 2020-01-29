package usecases

// GomsRepository interface that represents all the methods available to
// interact with premium-carousel-api microservice
type GomsRepository interface {
	GetHealthcheck() (string, error)
}

// UserBasicData is the structure that contains the basic user data
type UserBasicData struct {
	Name    string
	Phone   string
	Gender  string
	Country string
	Region  string
	Commune string
}

// UserRepository defines the methods that a User repository should have
type UserRepository interface {
	// GetUserData gets the user data based on his email
	GetUserData(email string) (UserBasicData, error)
}

// UserProfileRepository defines the methods that a User Profile repository should have
type UserProfileRepository interface {
	// GetUserData gets the user data based on his email
	GetUserProfileData(email string) (UserBasicData, error)
}

// SearchResponse object to recieve search response data type
type SearchResponse map[string][]struct {
	ListID        string `json:"listId"`
	Subject       string `json:"subject"`
	Price         string `json:"price"`
	UnitOfAccount string `json:"priceUf"`
	Currency      string `json:"currency"`
	Image         Image  `json:"image"`
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
