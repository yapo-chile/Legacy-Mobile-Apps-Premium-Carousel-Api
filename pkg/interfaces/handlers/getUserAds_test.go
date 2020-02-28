package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

func TestGetUserAdsHandlerInput(t *testing.T) {
	var h GetUserAdsHandler
	mMockInputRequest := &MockInputRequest{}
	mTargetRequest := &MockTargetRequest{}
	mMockInputRequest.On("Set",
		mock.AnythingOfType("*handlers.getUserAdsHandlerInput")).Return(mTargetRequest)
	mTargetRequest.On("FromPath").Return(mTargetRequest)
	input := h.Input(mMockInputRequest)
	var expected *getUserAdsHandlerInput
	assert.IsType(t, expected, input)
	mMockInputRequest.AssertExpectations(t)
	mTargetRequest.AssertExpectations(t)
}

type mockGetUserAdsInteractor struct {
	mock.Mock
}

func (m *mockGetUserAdsInteractor) GetUserAds(currentAdview domain.Ad) (domain.Ads, error) {
	args := m.Called(currentAdview)
	return args.Get(0).(domain.Ads), args.Error(1)
}

type mockGetAdInteractor struct {
	mock.Mock
}

func (m *mockGetAdInteractor) GetAd(listID string) (domain.Ad, error) {
	args := m.Called(listID)
	return args.Get(0).(domain.Ad), args.Error(1)
}

func TestGetUserAdsHandlerOK(t *testing.T) {
	mInteractor := &mockGetUserAdsInteractor{}
	mGetAdInteractor := &mockGetAdInteractor{}
	mGetAdInteractor.On("GetAd", mock.AnythingOfType("string")).
		Return(domain.Ad{ID: "123", UserID: "465"}, nil)
	mInteractor.On("GetUserAds", mock.AnythingOfType("domain.Ad")).
		Return(domain.Ads{{ID: "321", UserID: "465"}}, nil)
	h := GetUserAdsHandler{
		Interactor:      mInteractor,
		GetAdInteractor: mGetAdInteractor,
	}
	var input getUserAdsHandlerInput
	input.ListID = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getUserRequestOutput{
			Ads: []adsOutput{{ID: "321"}},
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mGetAdInteractor.AssertExpectations(t)
}
func TestGetUserAdsHandlerWithUF(t *testing.T) {
	mInteractor := &mockGetUserAdsInteractor{}
	mGetAdInteractor := &mockGetAdInteractor{}
	mGetAdInteractor.On("GetAd", mock.AnythingOfType("string")).
		Return(domain.Ad{ID: "123", UserID: "465"}, nil)
	mInteractor.On("GetUserAds", mock.AnythingOfType("domain.Ad")).
		Return(domain.Ads{{ID: "321", UserID: "465", Currency: "uf"}}, nil)
	h := GetUserAdsHandler{
		Interactor:          mInteractor,
		GetAdInteractor:     mGetAdInteractor,
		UnitOfAccountSymbol: "UF",
	}
	var input getUserAdsHandlerInput
	input.ListID = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getUserRequestOutput{
			Ads: []adsOutput{{ID: "321", Currency: "UF"}},
		},
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mGetAdInteractor.AssertExpectations(t)
}

func TestGetUserAdsHandlerNoAds(t *testing.T) {
	mInteractor := &mockGetUserAdsInteractor{}
	mGetAdInteractor := &mockGetAdInteractor{}
	mGetAdInteractor.On("GetAd", mock.AnythingOfType("string")).
		Return(domain.Ad{ID: "123", UserID: "465"}, nil)
	mInteractor.On("GetUserAds", mock.AnythingOfType("domain.Ad")).
		Return(domain.Ads{}, nil)
	h := GetUserAdsHandler{
		Interactor:      mInteractor,
		GetAdInteractor: mGetAdInteractor,
	}
	var input getUserAdsHandlerInput
	input.ListID = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mGetAdInteractor.AssertExpectations(t)
}

func TestGetUserAdsHandlerErrorGettingUserAds(t *testing.T) {
	mInteractor := &mockGetUserAdsInteractor{}
	mGetAdInteractor := &mockGetAdInteractor{}
	mGetAdInteractor.On("GetAd", mock.AnythingOfType("string")).
		Return(domain.Ad{ID: "123", UserID: "465"}, nil)
	mInteractor.On("GetUserAds", mock.AnythingOfType("domain.Ad")).
		Return(domain.Ads{}, fmt.Errorf("e"))
	h := GetUserAdsHandler{
		Interactor:      mInteractor,
		GetAdInteractor: mGetAdInteractor,
	}
	var input getUserAdsHandlerInput
	input.ListID = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mGetAdInteractor.AssertExpectations(t)
}

func TestGetUserAdsHandlerErrorGettingAd(t *testing.T) {
	mInteractor := &mockGetUserAdsInteractor{}
	mGetAdInteractor := &mockGetAdInteractor{}
	mGetAdInteractor.On("GetAd", mock.AnythingOfType("string")).
		Return(domain.Ad{ID: "123", UserID: "465"}, fmt.Errorf("e"))
	h := GetUserAdsHandler{
		Interactor:      mInteractor,
		GetAdInteractor: mGetAdInteractor,
	}
	var input getUserAdsHandlerInput
	input.ListID = "123"
	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mGetAdInteractor.AssertExpectations(t)
}

func TestGetUserAdsHandlerErrorBadInput(t *testing.T) {
	mInteractor := &mockGetUserAdsInteractor{}
	mGetAdInteractor := &mockGetAdInteractor{}
	h := GetUserAdsHandler{
		Interactor:      mInteractor,
		GetAdInteractor: mGetAdInteractor,
	}
	var input getUserAdsHandlerInput
	input.ListID = "123"
	getter := MakeMockInputGetter(&input, &goutils.Response{
		Code: http.StatusNoContent,
	})
	r := h.Execute(getter)
	expected := &goutils.Response{
		Code: http.StatusNoContent,
	}
	assert.Equal(t, expected, r)
	mInteractor.AssertExpectations(t)
	mGetAdInteractor.AssertExpectations(t)
}
