package usecases

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type mockPurchaseRepo struct {
	mock.Mock
}

func (m *mockPurchaseRepo) CreatePurchase(purchaseNumber, price int,
	purchaseType domain.PurchaseType) (domain.Purchase, error) {
	args := m.Called(purchaseNumber, price, purchaseType)
	return args.Get(0).(domain.Purchase), args.Error(1)
}

func (m *mockPurchaseRepo) AcceptPurchase(purchase domain.Purchase) (domain.Purchase, error) {
	args := m.Called(purchase)
	return args.Get(0).(domain.Purchase), args.Error(1)
}

type mockAddUserProductLogger struct {
	mock.Mock
}

func (m *mockAddUserProductLogger) LogErrorAddingProduct(UserID int, err error) {
	m.Called(UserID, err)
}

func (m *mockAddUserProductLogger) LogWarnSettingCache(UserID int, err error) {
	m.Called(UserID, err)
}

func TestAddProductOk(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("domain.Product"),
		mock.Anything).
		Return(nil)
	mPurchaseRepo.On("CreatePurchase",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType")).Return(domain.Purchase{}, nil)
	mProductRepo.On("GetUserActiveProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductType")).Return(domain.Product{},
		ErrProductNotFound)
	mProductRepo.On("CreateUserProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("domain.Purchase"),
		domain.PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("domain.ProductParams"),
	).Return(product, nil)
	mPurchaseRepo.On("AcceptPurchase",
		mock.AnythingOfType("domain.Purchase")).Return(domain.Purchase{}, nil)
	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestAddProductValidateError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0)

	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)
	mProductRepo.On("GetUserActiveProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductType")).Return(domain.Product{},
		nil)

	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestAddProductCreatePurchaseError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0)
	mPurchaseRepo.On("CreatePurchase",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType")).
		Return(domain.Purchase{}, fmt.Errorf("err"))
	mProductRepo.On("GetUserActiveProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductType")).Return(domain.Product{},
		ErrProductNotFound)

	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestAddProductAcceptPurchaseError(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0)
	mPurchaseRepo.On("CreatePurchase",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType")).Return(domain.Purchase{}, nil)
	mProductRepo.On("GetUserActiveProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductType")).Return(domain.Product{},
		ErrProductNotFound)
	mProductRepo.On("CreateUserProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("domain.Purchase"),
		domain.PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("domain.ProductParams"),
	).Return(product, nil)
	mPurchaseRepo.On("AcceptPurchase",
		mock.AnythingOfType("domain.Purchase")).
		Return(domain.Purchase{}, fmt.Errorf("err"))
	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestAddProductErrorAddingProduct(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0)
	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)
	mProductRepo.On("GetUserActiveProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductType")).Return(domain.Product{},
		ErrProductNotFound)
	mPurchaseRepo.On("CreatePurchase",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType")).Return(domain.Purchase{}, nil)
	mProductRepo.On("CreateUserProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("domain.Purchase"),
		domain.PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("domain.ProductParams"),
	).Return(product, fmt.Errorf("err"))
	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestAddProductOkErrorSettingCache(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("domain.Product"),
		mock.Anything).
		Return(fmt.Errorf("err"))
	mProductRepo.On("GetUserActiveProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.ProductType")).Return(domain.Product{},
		ErrProductNotFound)
	mPurchaseRepo.On("CreatePurchase",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType")).Return(domain.Purchase{}, nil)
	mProductRepo.On("CreateUserProduct",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("domain.Purchase"),
		domain.PremiumCarousel,
		mock.AnythingOfType("time.Time"),
		mock.AnythingOfType("domain.ProductParams"),
	).Return(product, nil)
	mPurchaseRepo.On("AcceptPurchase",
		mock.AnythingOfType("domain.Purchase")).Return(domain.Purchase{}, nil)
	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
