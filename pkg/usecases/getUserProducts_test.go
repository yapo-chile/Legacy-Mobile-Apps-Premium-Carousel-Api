package usecases

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockGetUserProductsLogger struct {
	mock.Mock
}

func (m *mockAddUserProductLogger) LogErrorGettingUserProducts(email string, err error) {
	m.Called(email, err)
}

func TestGetUserProductsOk(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeGetUserProductsInteractor(mProductRepo, mLogger)
	products := []Product{}
	mProductRepo.On("GetUserProducts",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
	).Return(products, 1, 1, nil)
	res, currentPage, totalPages, err := interactor.GetUserProducts("", 1)
	assert.NoError(t, err)
	assert.Equal(t, products, res)
	assert.Equal(t, 1, currentPage)
	assert.Equal(t, 1, totalPages)
	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserProductsError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeGetUserProductsInteractor(mProductRepo, mLogger)
	mLogger.On("LogErrorGettingUserProducts", mock.Anything, mock.Anything)
	mProductRepo.On("GetUserProducts",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
	).Return([]Product{}, 0, 0, fmt.Errorf("err"))
	_, _, _, err := interactor.GetUserProducts("", 1)
	assert.Error(t, err)

	mProductRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
