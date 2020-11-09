package usecases

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockExpireProductsLogger struct {
	mock.Mock
}

func (m *mockExpireProductsLogger) LogExpireProductsError(err error) {
	m.Called(err)
}

func (m *mockExpireProductsLogger) LogErrorSettingCache(err error) {
	m.Called(err)
}

func TestExpireProductsOk(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockExpireProductsLogger{}
	interactor := MakeExpireProductsInteractor(mProductRepo, mLogger)
	mProductRepo.On("ExpireProducts").Return(nil)
	err := interactor.ExpireProducts()
	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestExpireProductsRepoError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockExpireProductsLogger{}
	interactor := MakeExpireProductsInteractor(mProductRepo, mLogger)
	mProductRepo.On("ExpireProducts").Return(fmt.Errorf("err"))
	mLogger.On("LogExpireProductsError", mock.Anything)
	err := interactor.ExpireProducts()
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
}
