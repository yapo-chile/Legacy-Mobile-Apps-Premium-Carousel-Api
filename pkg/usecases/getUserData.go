package usecases

import (
	"fmt"
)

type GetUserDataInteractor interface {
	// GetUserData gets the user data based on his email
	GetUserData(pSHA1E string) (UserBasicData, error)
}

// getUserDataInteractor defines the interactor
type getUserDataInteractor struct {
	userProfileRepo UserProfileRepository
}

func MakeGetUserDataInteractor(userProfileRepo UserProfileRepository) GetUserDataInteractor {
	return &getUserDataInteractor{userProfileRepo: userProfileRepo}
}

// GetUser retrieves the basic data of a user given a mail
func (interactor *getUserDataInteractor) GetUserData(pSHA1Email string) (UserBasicData, error) {
	userProfile, err := interactor.userProfileRepo.GetUserProfileData(pSHA1Email)
	if err != nil {
		return userProfile, fmt.Errorf("error: cannot retrieve the user's profile: %+v", err)
	}
	return userProfile, nil
}
