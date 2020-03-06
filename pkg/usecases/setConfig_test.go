package usecases

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSetConfigLogger struct {
	mock.Mock
}

func (m *mockSetConfigLogger) LogErrorSettingConfig(userProductID int, err error) {
	m.Called(userProductID, err)
}

func (m *mockSetConfigLogger) LogWarnSettingCache(UserID string, err error) {
	m.Called(UserID, err)
}

func TestSetConfigOK(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetConfigLogger{}
	interactor := MakeSetConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mProductRepo.On("SetExpiration", mock.AnythingOfType("int"),
		mock.Anything).Return(nil)
	mProductRepo.On("SetConfig", mock.AnythingOfType("int"),
		mock.AnythingOfType("CpConfig")).
		Return(nil)
	mProductRepo.On("GetUserProductByID", mock.AnythingOfType("int")).
		Return(Product{}, nil)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("Product"),
		mock.Anything).
		Return(nil)
	err := interactor.SetConfig(1, CpConfig{}, time.Now().Add(time.Hour))
	assert.NoError(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}

func TestSetConfigErrorOnSetExpiration(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetConfigLogger{}
	interactor := MakeSetConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mProductRepo.On("SetExpiration", mock.AnythingOfType("int"),
		mock.Anything).Return(fmt.Errorf("err"))
	mLogger.On("LogErrorSettingConfig",
		mock.AnythingOfType("int"), mock.Anything)
	err := interactor.SetConfig(1, CpConfig{}, time.Now().Add(time.Hour))
	assert.Error(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}

func TestSetConfigOKErrorOnSetConfig(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetConfigLogger{}
	interactor := MakeSetConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mProductRepo.On("SetExpiration", mock.AnythingOfType("int"),
		mock.Anything).Return(nil)
	mProductRepo.On("SetConfig", mock.AnythingOfType("int"),
		mock.AnythingOfType("CpConfig")).
		Return(fmt.Errorf("err"))
	mLogger.On("LogErrorSettingConfig",
		mock.AnythingOfType("int"), mock.Anything)
	err := interactor.SetConfig(1, CpConfig{}, time.Now().Add(time.Hour))
	assert.Error(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}

func TestSetConfigOKErrorOnGetUserProductByID(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockSetConfigLogger{}
	interactor := MakeSetConfigInteractor(mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	mProductRepo.On("SetExpiration", mock.AnythingOfType("int"),
		mock.Anything).Return(nil)
	mProductRepo.On("SetConfig", mock.AnythingOfType("int"),
		mock.AnythingOfType("CpConfig")).
		Return(nil)
	mProductRepo.On("GetUserProductByID", mock.AnythingOfType("int")).
		Return(Product{}, fmt.Errorf("err"))
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("Product"),
		mock.Anything).
		Return(fmt.Errorf("err"))
	mLogger.On("LogWarnSettingCache",
		mock.AnythingOfType("string"), mock.Anything)
	err := interactor.SetConfig(1, CpConfig{}, time.Now().Add(time.Hour))
	assert.NoError(t, err)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mProductRepo.AssertExpectations(t)
}
