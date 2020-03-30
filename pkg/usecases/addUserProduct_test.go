package usecases

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAddUserProductLogger struct {
	mock.Mock
}

func (m *mockAddUserProductLogger) LogErrorAddingProduct(UserID string, err error) {
	m.Called(UserID, err)
}

func (m *mockAddUserProductLogger) LogWarnSettingCache(UserID string, err error) {
	m.Called(UserID, err)
}

func TestAddProductOk(t *testing.T) {
	product := Product{}
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mCacheRepo, mLogger, 0)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("usecases.Product"),
		mock.Anything).
		Return(nil)
	mProductRepo.On("AddUserProduct",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("usecases.ProductParams"),
	).Return(product, nil)
	err := interactor.AddUserProduct("", "", "", PremiumCarousel,
		time.Time{}, ProductParams{})
	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestAddProductErrorAddingProduct(t *testing.T) {
	product := Product{}
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mCacheRepo, mLogger, 0)
	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)
	mProductRepo.On("AddUserProduct",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("usecases.ProductParams"),
	).Return(product, fmt.Errorf("err"))
	err := interactor.AddUserProduct("", "", "", PremiumCarousel,
		time.Time{}, ProductParams{})
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestAddProductOkErrorSettingCache(t *testing.T) {
	product := Product{}
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mCacheRepo, mLogger, 0)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("usecases.Product"),
		mock.Anything).
		Return(fmt.Errorf("err"))
	mProductRepo.On("AddUserProduct",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("usecases.ProductParams"),
	).Return(product, nil)
	err := interactor.AddUserProduct("", "", "", PremiumCarousel,
		time.Time{}, ProductParams{})
	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
