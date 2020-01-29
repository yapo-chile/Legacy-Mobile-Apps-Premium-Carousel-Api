package repository

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestNewGomsRepository(t *testing.T) {
	m := MockHTTPHandler{}

	s := ""
	expected := GomsRepository{
		Handler: &m,
		Path:    s,
		TimeOut: 40,
	}

	result := NewGomsRepository(&m, 40, s)

	assert.Equal(t, &expected, result)
	m.AssertExpectations(t)
}

func TestGetOK(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}
	response := GomsResponse{
		Status: "OK",
	}
	jsonResponse, _ := json.Marshal(response)

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(string(jsonResponse), nil).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, "OK", result)
	assert.Nil(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(nil, errors.New("Error")).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, "", result)
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestPostParseError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return("", nil).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, result, "")
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetEmptyResponseError(t *testing.T) {
	mHandler := MockHTTPHandler{}
	mRequest := MockRequest{}

	gomsRepo := GomsRepository{
		Handler: &mHandler,
	}

	mRequest.On("SetMethod", "GET").Return(&mRequest).Once()
	mRequest.On("SetPath", "").Return(&mRequest).Once()
	mRequest.On("SetTimeOut", mock.AnythingOfType("int")).Return(&mRequest).Once()

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return("", nil).Once()

	result, err := gomsRepo.GetHealthcheck()

	assert.Equal(t, result, "")
	assert.Error(t, err)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}
