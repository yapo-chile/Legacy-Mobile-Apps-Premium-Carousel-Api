package repository

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"
)

func TestNewHttpAdRepo(t *testing.T) {
	m := mockHTTPHandler{}
	expected := HTTPAdRepo{
		Handler: &m,
	}
	result := NewHTTPAdRepo(&m, "")
	assert.Equal(t, &expected, result)
	m.AssertExpectations(t)
}

func TestGetAd(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}
	adRepo := HTTPAdRepo{
		Handler: &mHandler,
	}
	defaultStr := `{"1":[{"listId":"1"}]}`
	expected := usecases.SearchResponse{}
	mRequest.On("SetMethod", "POST").Return(&mRequest)
	mRequest.On("SetPath", "").Return(&mRequest)
	mRequest.On("SetBody",
		mock.AnythingOfType("usecases.SearchInput")).Return(&mRequest)
	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(defaultStr, nil).Once()
	result, err := adRepo.Get(usecases.SearchInput{})

	json.Unmarshal([]byte(defaultStr), &expected)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetAdError(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}
	adRepo := HTTPAdRepo{
		Handler: &mHandler,
	}
	defaultStr := `{"1":[{"listId":"1"}]}`
	expected := usecases.SearchResponse{}
	mRequest.On("SetMethod", "POST").Return(&mRequest)
	mRequest.On("SetPath", "").Return(&mRequest)
	mRequest.On("SetBody",
		mock.AnythingOfType("usecases.SearchInput")).Return(&mRequest)
	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(defaultStr, fmt.Errorf("error")).Once()
	result, err := adRepo.Get(usecases.SearchInput{})

	assert.Error(t, err)
	assert.Equal(t, expected, result)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetAderrorBadJSON(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}
	adRepo := HTTPAdRepo{
		Handler: &mHandler,
	}
	defaultStr := `{"1":[{"listId":"1"]}`
	expected := usecases.SearchResponse{}
	mRequest.On("SetMethod", "POST").Return(&mRequest)
	mRequest.On("SetPath", "").Return(&mRequest)
	mRequest.On("SetBody",
		mock.AnythingOfType("usecases.SearchInput")).Return(&mRequest)

	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(defaultStr, nil).Once()
	result, err := adRepo.Get(usecases.SearchInput{})

	assert.Error(t, err)
	assert.Equal(t, expected, result)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}

func TestGetAdErroAdResultZero(t *testing.T) {
	mHandler := mockHTTPHandler{}
	mRequest := mockRequest{}
	adRepo := HTTPAdRepo{
		Handler: &mHandler,
	}
	defaultStr := `{}`
	expected := usecases.SearchResponse{}
	mRequest.On("SetMethod", "POST").Return(&mRequest)
	mRequest.On("SetPath", "").Return(&mRequest)
	mRequest.On("SetBody",
		mock.AnythingOfType("usecases.SearchInput")).Return(&mRequest)
	mHandler.On("NewRequest").Return(&mRequest, nil).Once()
	mHandler.On("Send", &mRequest).Return(defaultStr, nil).Once()
	result, err := adRepo.Get(usecases.SearchInput{})

	json.Unmarshal([]byte(defaultStr), &expected)
	assert.Error(t, err)
	assert.Equal(t, expected, result)
	mHandler.AssertExpectations(t)
	mRequest.AssertExpectations(t)
}
