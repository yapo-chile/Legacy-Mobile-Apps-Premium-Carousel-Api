package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

func TestMakePurchaseRepositoryOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	repo := MakePurchaseRepository(mockDB)
	assert.Equal(t, &purchaseRepo{
		handler: mockDB,
	}, repo)
	mockDB.AssertExpectations(t)
}

func TestCreatePurchaseOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	testTime := time.Now()
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil)
	mResult.On("Close").Return(nil)
	mResult.On("Next").Return(true).Once()
	mResult.On("Scan", mock.Anything).
		Return([]interface{}{123, testTime, domain.PendingPurchase})
	repo := MakePurchaseRepository(mockDB)
	result, err := repo.CreatePurchase(10, 100, domain.AdminPurchase)
	assert.NoError(t, err)
	expected := domain.Purchase{
		ID:        123,
		Number:    10,
		Price:     100,
		Type:      domain.AdminPurchase,
		Status:    domain.PendingPurchase,
		CreatedAt: testTime,
	}
	assert.Equal(t, expected, result)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestCreatePurchaseQueryError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("err"))
	repo := MakePurchaseRepository(mockDB)
	_, err := repo.CreatePurchase(10, 100, domain.AdminPurchase)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestCreatePurchaseNextError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil)
	mResult.On("Close").Return(nil)
	mResult.On("Next").Return(false).Once()
	repo := MakePurchaseRepository(mockDB)
	_, err := repo.CreatePurchase(10, 100, domain.AdminPurchase)
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestAcceptPurchaseOK(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	testTime := time.Now()
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, nil)
	mResult.On("Close").Return(nil)
	repo := MakePurchaseRepository(mockDB)
	prevPurchase := domain.Purchase{
		ID:        123,
		Number:    10,
		Price:     100,
		Type:      domain.AdminPurchase,
		Status:    domain.PendingPurchase,
		CreatedAt: testTime,
	}
	newPurchase, err := repo.AcceptPurchase(prevPurchase)
	expected := prevPurchase
	expected.Status = domain.AcceptedPurchase
	assert.NoError(t, err)
	assert.Equal(t, expected, newPurchase)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestAcceptPurchaseError(t *testing.T) {
	mockDB := &dbHandlerMock{}
	mResult := &mockResult{}
	mockDB.On("Query",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]interface {}"),
	).Return(mResult, fmt.Errorf("err"))
	repo := MakePurchaseRepository(mockDB)
	_, err := repo.AcceptPurchase(domain.Purchase{})
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mResult.AssertExpectations(t)
}
