package repository

import (
	"encoding/json"
	"fmt"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

const (
	errorNoUserDataFound string = "there was no user data found using the email: %s. %+v"
	errorUnmarshal       string = "there was an error parsing the user data %s"
)

// UserProfileRepository wrapper struct for the RedisHandler
type UserProfileRepository struct {
	Handler HTTPHandler
	Path    string
}

// MakeUserProfileRepository constructor
func MakeUserProfileRepository(handler HTTPHandler, path string) usecases.UserProfileRepository {
	return &UserProfileRepository{
		Handler: handler,
		Path:    path,
	}
}

// GetUserProfileData makes a http request to profile service
// to get the user profile data
// it sends the SHA1 representation of the provided email
func (repo *UserProfileRepository) GetUserProfileData(SHA1Email string) (usecases.UserBasicData, error) {
	request := repo.Handler.NewRequest().
		SetMethod("GET").SetPath(fmt.Sprintf(repo.Path, SHA1Email))

	JSONResp, err := repo.Handler.Send(request)
	if err == nil && JSONResp != "" {
		resp := fmt.Sprintf("%s", JSONResp)
		var userData map[string]usecases.UserBasicData

		err := json.Unmarshal([]byte(resp), &userData)
		if err != nil {
			return usecases.UserBasicData{}, fmt.Errorf(errorUnmarshal, SHA1Email, err)
		}

		val, ok := userData[SHA1Email]
		if !ok {
			return usecases.UserBasicData{}, fmt.Errorf(errorNoUserDataFound, SHA1Email)
		}

		return val, err
	}

	return usecases.UserBasicData{}, fmt.Errorf(errorNoUserDataFound, SHA1Email)
}
