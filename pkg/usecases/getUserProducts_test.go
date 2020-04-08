package usecases

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type mockGetUserProductsLogger struct {
	mock.Mock
}

func (m *mockGetUserProductsLogger) LogErrorGettingUserProducts(err error) {
	m.Called(err)
}

func (m *mockGetUserProductsLogger) LogErrorGettingUserProductsByEmail(email string, err error) {
	m.Called(email, err)
}

func TestGetUserProductsByEmailOk(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockGetUserProductsLogger{}
	interactor := MakeGetUserProductsInteractor(mProductRepo, mLogger)
	products := []domain.Product{}
	mProductRepo.On("GetUserProductsByEmail",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
	).Return(products, 1, 1, nil)
	res, currentPage,
		totalPages, err := interactor.GetUserProducts("test@test.cl", 1)
	assert.NoError(t, err)
	assert.Equal(t, products, res)
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 1, totalPages)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsByEmailError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockGetUserProductsLogger{}
	interactor := MakeGetUserProductsInteractor(mProductRepo, mLogger)
	mLogger.On("LogErrorGettingUserProductsByEmail",
		mock.Anything, mock.Anything)
	mProductRepo.On("GetUserProductsByEmail",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
	).Return([]domain.Product{}, 0, 0, fmt.Errorf("err"))
	_, _, _, err := interactor.GetUserProducts("test@test.cl", 1)
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsOk(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockGetUserProductsLogger{}
	interactor := MakeGetUserProductsInteractor(mProductRepo, mLogger)
	products := []domain.Product{}
	mProductRepo.On("GetUserProducts",
		mock.AnythingOfType("int"),
	).Return(products, 1, 1, nil)
	res, currentPage,
		totalPages, err := interactor.GetUserProducts("", 1)
	assert.NoError(t, err)
	assert.Equal(t, products, res)
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 1, totalPages)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockGetUserProductsLogger{}
	interactor := MakeGetUserProductsInteractor(mProductRepo, mLogger)
	mLogger.On("LogErrorGettingUserProducts",
		mock.Anything, mock.Anything)
	mProductRepo.On("GetUserProducts",
		mock.AnythingOfType("int"),
	).Return([]domain.Product{}, 0, 0, fmt.Errorf("err"))
	_, _, _, err := interactor.GetUserProducts("", 1)
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
