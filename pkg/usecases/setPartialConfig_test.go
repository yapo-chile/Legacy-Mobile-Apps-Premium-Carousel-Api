package usecases

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSetPartialConfigLogger struct {
	mock.Mock
}

func (m *mockSetPartialConfigLogger) LogErrorSettingPartialConfig(userProductID int, err error) {
	m.Called(userProductID, err)
}

func (m *mockSetPartialConfigLogger) LogWarnSettingCache(UserID string, err error) {
	m.Called(UserID, err)
}

func TestSetPartialConfigOK(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetPartialConfigLogger{}
	interactor := MakeSetPartialConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mProductRepo.On("SetPartialConfig", mock.AnythingOfType("int"),
		mock.Anything).Return(nil)
	mProductRepo.On("GetUserProductByID", mock.AnythingOfType("int")).
		Return(Product{}, nil)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("Product"),
		mock.Anything).
		Return(nil)
	err := interactor.SetPartialConfig(1, map[string]interface{}{})
	assert.NoError(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}

func TestSetPartialConfigErrorSettingPartialConfig(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetPartialConfigLogger{}
	interactor := MakeSetPartialConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mProductRepo.On("SetPartialConfig", mock.AnythingOfType("int"),
		mock.Anything).Return(fmt.Errorf("err"))
	mLogger.On("LogErrorSettingPartialConfig", mock.AnythingOfType("int"),
		mock.Anything)
	err := interactor.SetPartialConfig(1, map[string]interface{}{})
	assert.Error(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}

func TestSetPartialConfigErrorGettingProduct(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetPartialConfigLogger{}
	interactor := MakeSetPartialConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mLogger.On("LogWarnSettingCache", mock.AnythingOfType("string"),
		mock.Anything)
	mProductRepo.On("SetPartialConfig", mock.AnythingOfType("int"),
		mock.Anything).Return(nil)
	mProductRepo.On("GetUserProductByID", mock.AnythingOfType("int")).
		Return(Product{}, fmt.Errorf("err"))
	err := interactor.SetPartialConfig(1, map[string]interface{}{})
	assert.Error(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}

func TestSetPartialConfigErrorSettingCache(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetPartialConfigLogger{}
	interactor := MakeSetPartialConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mProductRepo.On("SetPartialConfig", mock.AnythingOfType("int"),
		mock.Anything).Return(nil)
	mProductRepo.On("GetUserProductByID", mock.AnythingOfType("int")).
		Return(Product{}, nil)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("Product"),
		mock.Anything).
		Return(fmt.Errorf("err"))
	mLogger.On("LogWarnSettingCache", mock.AnythingOfType("string"),
		mock.Anything)
	err := interactor.SetPartialConfig(1, map[string]interface{}{})
	assert.NoError(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}
