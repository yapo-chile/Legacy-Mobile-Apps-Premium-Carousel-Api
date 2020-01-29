package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

func TestNewUserDataRepository(t *testing.T) {
	mHandler := MockHTTPHandler{}
	expected := &UserProfileRepository{
		Handler: &mHandler,
	}
	repo := NewUserProfileRepository(&mHandler, "")
	assert.Equal(t, expected, repo)
	mHandler.AssertExpectations(t)
}

func TestUserDataRepositoryGetUserDataOK(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)

	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return(`{"9fbcd51ef0a2d6293730c6a60afee8c807677fb5":{"uuid":"edgar"}}`, nil)

	expected := usecases.UserBasicData{Name: ""}

	repo := UserProfileRepository{
		Handler: &mHandler,
	}
	resp, err := repo.GetUserProfileData("edgar@gmail.com")
	assert.Equal(t, expected, resp)
	assert.NoError(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestUserDataRepositoryGetUserDataError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)

	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return(`{"notasha1email":{}}`, nil)

	expected := usecases.UserBasicData{}

	repo := UserProfileRepository{
		Handler: &mHandler,
	}
	resp, err := repo.GetUserProfileData("edgar@gmail.com")
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestUserDataRepositoryGetUserDataRequestError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)

	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return("", fmt.Errorf("error"))

	expected := usecases.UserBasicData{}

	repo := UserProfileRepository{
		Handler: &mHandler,
	}
	resp, err := repo.GetUserProfileData("edgar@gmail.com")
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestUserDataRepositoryUnmarshalError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	mRequest.On("SetPath", mock.AnythingOfType("string")).Return(&mRequest)
	mRequest.On("SetMethod", "GET").Return(&mRequest)

	mHandler.On("NewRequest").Return(&mRequest)
	mHandler.On("Send", &mRequest).Return(`{"9fbcd51ef0a2d6293730c6a60afee8c807677fb5":"uuid":"edgar"}}`, nil).Once()

	expected := usecases.UserBasicData{}

	repo := UserProfileRepository{
		Handler: &mHandler,
	}
	resp, err := repo.GetUserProfileData("edgar@gmail.com")
	assert.Equal(t, expected, resp)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}
