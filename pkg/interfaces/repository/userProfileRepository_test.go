package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

type MockHTTPHandler struct {
	mock.Mock
}

func (m *MockHTTPHandler) Send(request HTTPRequest) (interface{}, error) {
	args := m.Called(request)
	return args.Get(0), args.Error(1)
}

func (m *MockHTTPHandler) NewRequest() HTTPRequest {
	args := m.Called()
	return args.Get(0).(HTTPRequest)
}

type MockRequest struct {
	mock.Mock
}

func (m *MockRequest) SetMethod(method string) HTTPRequest {
	args := m.Called(method)
	return args.Get(0).(HTTPRequest)
}

func (m *MockRequest) GetMethod() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRequest) SetPath(path string) HTTPRequest {
	args := m.Called(path)
	return args.Get(0).(HTTPRequest)
}

func (m *MockRequest) GetPath() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockRequest) SetHeaders(headers map[string]string) HTTPRequest {
	args := m.Called(headers)
	return args.Get(0).(HTTPRequest)
}

func (m *MockRequest) GetHeaders() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *MockRequest) SetBody(body interface{}) HTTPRequest {
	args := m.Called(body)
	return args.Get(0).(HTTPRequest)
}

func (m *MockRequest) GetBody() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockRequest) SetQueryParams(queryParams map[string]string) HTTPRequest {
	args := m.Called(queryParams)
	return args.Get(0).(HTTPRequest)
}

func (m *MockRequest) GetQueryParams() map[string][]string {
	args := m.Called()
	return args.Get(0).(map[string][]string)
}

func (m *MockRequest) GetTimeOut() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func (m *MockRequest) SetTimeOut(t int) HTTPRequest {
	args := m.Called(t)
	return args.Get(0).(HTTPRequest)
}

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
