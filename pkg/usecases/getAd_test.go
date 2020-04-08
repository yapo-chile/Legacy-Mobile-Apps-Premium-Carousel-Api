package usecases

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type mockGetAdLogger struct {
	mock.Mock
}

func (m *mockGetAdLogger) LogWarnGettingCache(listID string, err error) {
	m.Called(listID, err)
}

func (m *mockGetAdLogger) LogWarnSettingCache(listID string, err error) {
	m.Called(listID, err)
}

func (m *mockGetAdLogger) LogErrorGettingAd(listID string, err error) {
	m.Called(listID, err)
}

func TestGetAdOkWithoutCache(t *testing.T) {
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockGetAdLogger{}
	interactor := MakeGetAdInteractor(mAdRepo, mCacheRepo, mLogger, 0)
	tAd := domain.Ad{ID: "1", Subject: "Mi auto", UserID: 123}
	mLogger.On("LogWarnGettingCache", mock.Anything, mock.Anything)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), mock.Anything).
		Return([]byte{}, fmt.Errorf("cache not found"))
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		MinifiedAdDataType,
		mock.AnythingOfType("domain.Ad"),
		mock.Anything).
		Return(nil)
	mAdRepo.On("GetAd", mock.AnythingOfType("string")).Return(tAd, nil)
	ads, err := interactor.GetAd("1")
	expected := tAd
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
}

func TestGetAdOkWithCache(t *testing.T) {
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockGetAdLogger{}
	interactor := MakeGetAdInteractor(mAdRepo, mCacheRepo, mLogger, 0)
	tAd := domain.Ad{ID: "1", Subject: "Mi auto", UserID: 123}
	tAdBytes, _ := json.Marshal(tAd)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), mock.Anything).
		Return(tAdBytes, nil)
	ads, err := interactor.GetAd("1")
	expected := tAd
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
}

func TestGetAdErrorGettingAd(t *testing.T) {
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockGetAdLogger{}
	interactor := MakeGetAdInteractor(mAdRepo, mCacheRepo, mLogger, 0)
	tAd := domain.Ad{ID: "1", Subject: "Mi auto", UserID: 123}
	mLogger.On("LogWarnGettingCache", mock.Anything, mock.Anything)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), mock.Anything).
		Return([]byte{}, fmt.Errorf("cache not found"))
	mLogger.On("LogErrorGettingAd", mock.Anything, mock.Anything)
	mAdRepo.On("GetAd", mock.AnythingOfType("string")).Return(tAd, fmt.Errorf("err"))
	_, err := interactor.GetAd("1")
	assert.Error(t, err)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
}

func TestGetAdErrorSettingCache(t *testing.T) {
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockGetAdLogger{}
	interactor := MakeGetAdInteractor(mAdRepo, mCacheRepo, mLogger, 0)
	tAd := domain.Ad{ID: "1", Subject: "Mi auto", UserID: 123}
	mLogger.On("LogWarnGettingCache", mock.Anything, mock.Anything)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), mock.Anything).
		Return([]byte{}, fmt.Errorf("cache not found"))
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		MinifiedAdDataType,
		mock.AnythingOfType("domain.Ad"),
		mock.Anything).
		Return(fmt.Errorf("err"))
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)
	mAdRepo.On("GetAd", mock.AnythingOfType("string")).Return(tAd, nil)
	ads, err := interactor.GetAd("1")
	expected := tAd
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
}
