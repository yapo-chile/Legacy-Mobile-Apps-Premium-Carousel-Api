package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

type mockHTTPHandler struct {
	mock.Mock
}

func (m *mockHTTPHandler) Send(request HTTPRequest) (interface{}, error) {
	args := m.Called(request)
	return args.Get(0), args.Error(1)
}

func (m *mockHTTPHandler) NewRequest() HTTPRequest {
	args := m.Called()
	return args.Get(0).(HTTPRequest)
}

type mockRequest struct {
	mock.Mock
}

func (m *mockRequest) SetMethod(method string) HTTPRequest {
	args := m.Called(method)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetMethod() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRequest) SetPath(path string) HTTPRequest {
	args := m.Called(path)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRequest) SetHeaders(headers map[string]string) HTTPRequest {
	args := m.Called(headers)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetHeaders() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *mockRequest) SetBody(body interface{}) HTTPRequest {
	args := m.Called(body)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetBody() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *mockRequest) SetQueryParams(queryParams map[string]string) HTTPRequest {
	args := m.Called(queryParams)
	return args.Get(0).(HTTPRequest)
}

func (m *mockRequest) GetQueryParams() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *mockRequest) GetTimeOut() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *mockRequest) SetTimeOut(t int) HTTPRequest {
	args := m.Called(t)
	return args.Get(0).(HTTPRequest)
}

func TestNewUserDataRepository(t *testing.T) {
	mHandler := mockHTTPHandler{}
	expected := &UserProfileRepository{
		Handler: &mHandler,
	}
	repo := NewUserProfileRepository(&mHandler, "")
	assert.Equal(t, expected, repo)
	mHandler.AssertExpectations(t)
}

func TestUserDataRepositoryGetUserDataOK(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

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
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

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
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

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
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}

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
