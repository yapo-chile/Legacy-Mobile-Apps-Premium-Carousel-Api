package usecases

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type mockConfigRepo struct {
	mock.Mock
}

func (m *mockConfigRepo) GetConfig(userID string) (CpConfig, error) {
	args := m.Called(userID)
	return args.Get(0).(CpConfig), args.Error(1)
}

type mockAdRepo struct {
	mock.Mock
}

func (m *mockAdRepo) GetUserAds(userID string,
	cpConfig CpConfig) (domain.Ads, error) {
	args := m.Called(userID, cpConfig)
	return args.Get(0).(domain.Ads), args.Error(1)
}

func (m *mockAdRepo) GetAd(listID string) (domain.Ad, error) {
	args := m.Called(listID)
	return args.Get(0).(domain.Ad), args.Error(1)
}

func TestGetUserAdsOk(t *testing.T) {
	mConfigRepo := &mockConfigRepo{}
	mAdRepo := &mockAdRepo{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mConfigRepo)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}
	tAds := domain.Ads{
		{ID: "1", Subject: "Mi auto", UserID: "123"},
		{ID: "2", Subject: "Mi auto 2", UserID: "123"},
	}
	mConfigRepo.On("GetConfig", mock.AnythingOfType("string")).
		Return(cpConfig, nil)
	mAdRepo.On("GetUserAds", mock.AnythingOfType("string"),
		mock.AnythingOfType("CpConfig")).Return(tAds, nil)
	ads, err := interactor.GetUserAds("123",
		"test_excluded_list_id_1234", "excluded_list_id_2123")
	expected := tAds
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
	mConfigRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
}

func TestGetUserAdsErrorGettingConfig(t *testing.T) {
	mConfigRepo := &mockConfigRepo{}
	mAdRepo := &mockAdRepo{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mConfigRepo)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}
	mConfigRepo.On("GetConfig", mock.AnythingOfType("string")).
		Return(cpConfig, fmt.Errorf("e"))

	_, err := interactor.GetUserAds("123",
		"test_excluded_list_id_1234", "excluded_list_id_2123")
	assert.Error(t, err)
	mConfigRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
}

func TestGetUserAdsErrorGettingAds(t *testing.T) {
	mConfigRepo := &mockConfigRepo{}
	mAdRepo := &mockAdRepo{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mConfigRepo)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}
	mConfigRepo.On("GetConfig", mock.AnythingOfType("string")).
		Return(cpConfig, nil)
	mAdRepo.On("GetUserAds", mock.AnythingOfType("string"),
		mock.AnythingOfType("CpConfig")).Return(domain.Ads{}, fmt.Errorf("e"))
	_, err := interactor.GetUserAds("123",
		"test_excluded_list_id_1234", "excluded_list_id_2123")
	assert.Error(t, err)
	mConfigRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
}
