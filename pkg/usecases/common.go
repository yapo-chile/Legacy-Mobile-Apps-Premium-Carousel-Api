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
