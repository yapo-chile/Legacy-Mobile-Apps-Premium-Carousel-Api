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

func (m *mockAddUserProductLogger) LogErrorAddingProduct(userID int, err error) {
	m.Called(userID, err)
}

func (m *mockAddUserProductLogger) LogWarnSettingCache(userID int, err error) {
	m.Called(userID, err)
}

func (m *mockAddUserProductLogger) LogWarnPushingEvent(productID int, err error) {
	m.Called(productID, err)
}

type mockBackendEventRepo struct {
	mock.Mock
}

func (m *mockBackendEventRepo) PushSoldProduct(product domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func TestAddProductOk(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	mBackendEventRepo := &mockBackendEventRepo{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0, mBackendEventRepo, true)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("domain.Product"),
		mock.Anything).
		Return(nil)
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
	mBackendEventRepo.On("PushSoldProduct",
		mock.AnythingOfType("domain.Product")).Return(nil)
	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mBackendEventRepo.AssertExpectations(t)
}

func TestAddProductCreatePurchaseError(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	mBackendEventRepo := &mockBackendEventRepo{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0, mBackendEventRepo, false)
	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)
	mPurchaseRepo.On("CreatePurchase",
		mock.AnythingOfType("int"),
		mock.AnythingOfType("int"),
		mock.AnythingOfType("domain.PurchaseType")).
		Return(domain.Purchase{}, fmt.Errorf("err"))

	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mBackendEventRepo.AssertExpectations(t)
}

func TestAddProductAcceptPurchaseError(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	mBackendEventRepo := &mockBackendEventRepo{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0, mBackendEventRepo, false)
	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)

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
		mock.AnythingOfType("domain.Purchase")).
		Return(domain.Purchase{}, fmt.Errorf("err"))
	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mBackendEventRepo.AssertExpectations(t)
}

func TestAddProductErrorAddingProduct(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	mBackendEventRepo := &mockBackendEventRepo{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0, mBackendEventRepo, false)
	mLogger.On("LogErrorAddingProduct", mock.Anything, mock.Anything)
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
	mBackendEventRepo := &mockBackendEventRepo{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0, mBackendEventRepo, false)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("domain.Product"),
		mock.Anything).
		Return(fmt.Errorf("err"))
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
	mBackendEventRepo.AssertExpectations(t)
}

func TestAddProductOkBackendEventError(t *testing.T) {
	product := domain.Product{}
	mProductRepo := &mockProductRepo{}
	mPurchaseRepo := &mockPurchaseRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockAddUserProductLogger{}
	mBackendEventRepo := &mockBackendEventRepo{}
	interactor := MakeAddUserProductInteractor(mProductRepo, mPurchaseRepo,
		mCacheRepo, mLogger, 0, mBackendEventRepo, true)
	mCacheRepo.On("SetCache", mock.AnythingOfType("string"),
		ProductCacheType,
		mock.AnythingOfType("domain.Product"),
		mock.Anything).
		Return(nil)
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
	mBackendEventRepo.On("PushSoldProduct",
		mock.AnythingOfType("domain.Product")).Return(fmt.Errorf("err"))
	mLogger.On("LogWarnPushingEvent", mock.Anything, mock.Anything)

	err := interactor.AddUserProduct(0, "", 0, 0, domain.AdminPurchase,
		domain.PremiumCarousel, time.Time{}, domain.ProductParams{})
	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mPurchaseRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
	mBackendEventRepo.AssertExpectations(t)
}
