package usecases

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
)

type mockProductRepo struct {
	mock.Mock
}

func (m *mockProductRepo) GetUserActiveProduct(userID int,
	productType domain.ProductType) (domain.Product, error) {
	args := m.Called(userID, productType)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *mockProductRepo) GetUserProducts(page int) ([]domain.Product, int, int, error) {
	args := m.Called(page)
	return args.Get(0).([]domain.Product), args.Int(1), args.Int(2), args.Error(3)
}

func (m *mockProductRepo) GetReport(start, end time.Time) ([]domain.Product, error) {
	args := m.Called(start, end)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *mockProductRepo) GetUserProductsByEmail(email string, page int) ([]domain.Product, int, int, error) {
	args := m.Called(email, page)
	return args.Get(0).([]domain.Product), args.Int(1), args.Int(2), args.Error(3)
}

func (m *mockProductRepo) CreateUserProduct(userID int, email string, purchase domain.Purchase,
	productType domain.ProductType, expiredAt time.Time,
	config domain.ProductParams) (domain.Product, error) {
	args := m.Called(userID, email, purchase, productType, expiredAt, config)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *mockProductRepo) GetUserProductsTotal() (total int) {
	args := m.Called()
	return args.Int(0)
}

func (m *mockProductRepo) GetUserProductsTotalByEmail(email string) (total int) {
	args := m.Called(email)
	return args.Int(0)
}

func (m *mockProductRepo) GetUserProductByID(userProductID int) (domain.Product, error) {
	args := m.Called(userProductID)
	return args.Get(0).(domain.Product), args.Error(1)
}

func (m *mockProductRepo) SetConfig(userProductID int, config domain.ProductParams) error {
	args := m.Called(userProductID, config)
	return args.Error(0)
}

func (m *mockProductRepo) SetPartialConfig(userProductID int,
	config map[string]interface{}) error {
	args := m.Called(userProductID, config)
	return args.Error(0)
}

func (m *mockProductRepo) SetExpiration(userProductID int,
	expiredAt time.Time) error {
	args := m.Called(userProductID, expiredAt)
	return args.Error(0)
}

func (m *mockProductRepo) SetStatus(userProductID int,
	status domain.ProductStatus) error {
	args := m.Called(userProductID, status)
	return args.Error(0)
}

func (m *mockProductRepo) ExpireProducts() error {
	args := m.Called()
	return args.Error(0)
}

type mockAdRepo struct {
	mock.Mock
}

func (m *mockAdRepo) GetUserAds(userID int,
	productParams domain.ProductParams) (domain.Ads, error) {
	args := m.Called(userID, productParams)
	return args.Get(0).(domain.Ads), args.Error(1)
}

func (m *mockAdRepo) GetAd(listID string) (domain.Ad, error) {
	args := m.Called(listID)
	return args.Get(0).(domain.Ad), args.Error(1)
}

type mockCacheRepo struct {
	mock.Mock
}

func (m *mockCacheRepo) SetCache(key string, typ CacheType, data interface{},
	expiration time.Duration) error {
	args := m.Called(key, typ, data, expiration)
	return args.Error(0)
}

func (m *mockCacheRepo) GetCache(key string, typ CacheType) ([]byte, error) {
	args := m.Called(key, typ)
	return args.Get(0).([]byte), args.Error(1)
}

type mockgetUserAdsLogger struct {
	mock.Mock
}

func (m *mockgetUserAdsLogger) LogWarnGettingCache(userID int, err error) {
	m.Called(userID, err)
}

func (m *mockgetUserAdsLogger) LogWarnSettingCache(userID int, err error) {
	m.Called(userID, err)
}

func (m *mockgetUserAdsLogger) LogInfoActiveProductNotFound(userID int, product domain.Product) {
	m.Called(userID, product)
}

func (m *mockgetUserAdsLogger) LogInfoProductExpired(userID int, product domain.Product) {
	m.Called(userID, product)
}

func (m *mockgetUserAdsLogger) LogErrorGettingUserAdsData(userID int, err error) {
	m.Called(userID, err)
}

func (m *mockgetUserAdsLogger) LogNotEnoughAds(userID int) {
	m.Called(userID)
}

func TestGetUserAdsOkWithoutCache(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour, 2)
	productParams := domain.ProductParams{Categories: []int{2020}, Limit: 2, PriceRange: 200}
	tAds := domain.Ads{
		{ID: "1", Subject: "Mi auto", UserID: 123},
		{ID: "2", Subject: "Mi auto 2", UserID: 123},
	}
	testTime := time.Now().Add(time.Hour * 24)
	product := domain.Product{Config: productParams, UserID: 123,
		ExpiredAt: testTime, Status: domain.ActiveProduct}
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return([]byte{}, fmt.Errorf("cache not found"))
	mLogger.On("LogWarnGettingCache", mock.Anything, mock.Anything)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)
	mProductRepo.On("GetUserActiveProduct", mock.AnythingOfType("int"),
		domain.PremiumCarousel).Return(product, nil)
	mAdRepo.On("GetUserAds", mock.AnythingOfType("int"),
		mock.AnythingOfType("ProductParams")).Return(tAds, nil)
	mCacheRepo.On("SetCache", "user:123:PREMIUM_CAROUSEL",
		ProductCacheType,
		product,
		time.Hour).
		Return(fmt.Errorf("error setting cache"))
	ads, err := interactor.GetUserAds(domain.Ad{UserID: 123})
	expected := tAds
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserAdsOkWithoutCacheAndInactiveProduct(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour, 2)

	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return([]byte{}, fmt.Errorf("cache not found"))
	mLogger.On("LogWarnGettingCache", mock.Anything, mock.Anything)
	mLogger.On("LogInfoActiveProductNotFound", mock.Anything, mock.Anything)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)
	mProductRepo.On("GetUserActiveProduct", mock.AnythingOfType("int"),
		domain.PremiumCarousel).Return(domain.Product{}, fmt.Errorf("err"))

	product := domain.Product{UserID: 123, Status: domain.InactiveProduct}

	mCacheRepo.On("SetCache", "user:123:PREMIUM_CAROUSEL",
		ProductCacheType,
		product,
		time.Hour).
		Return(fmt.Errorf("error setting cache"))
	_, err := interactor.GetUserAds(domain.Ad{UserID: 123})

	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserAdsOkWithCache(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour, 2)
	productParams := domain.ProductParams{Categories: []int{2020}, Limit: 2}
	tAds := domain.Ads{
		{ID: "1", Subject: "Mi auto", UserID: 123},
		{ID: "2", Subject: "Mi auto 2", UserID: 123},
	}
	testTime := time.Now().Add(time.Hour * 24)
	product := domain.Product{Config: productParams,
		ExpiredAt: testTime, Status: domain.ActiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mAdRepo.On("GetUserAds", mock.AnythingOfType("int"),
		mock.AnythingOfType("ProductParams")).Return(tAds, nil)
	ads, err := interactor.GetUserAds(domain.Ad{UserID: 123})
	expected := tAds
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserAdsErrorProductInactive(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour, 2)
	productParams := domain.ProductParams{Categories: []int{2020}, Limit: 2}
	testTime := time.Now().Add(time.Hour * 24)
	product := domain.Product{Config: productParams,
		ExpiredAt: testTime, Status: domain.InactiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mLogger.On("LogInfoActiveProductNotFound", mock.Anything, mock.Anything)

	interactor.GetUserAds(domain.Ad{UserID: 123})
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserAdsErrorProductExpired(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour, 2)
	productParams := domain.ProductParams{Categories: []int{2020}, Limit: 2}
	testTime := time.Now().Add(time.Hour * -24)
	product := domain.Product{Config: productParams,
		ExpiredAt: testTime, UserID: 123, Status: domain.ActiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mLogger.On("LogInfoProductExpired", mock.Anything, mock.Anything)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)

	product.Status = domain.ExpiredProduct
	mCacheRepo.On("SetCache", "user:123:PREMIUM_CAROUSEL",
		ProductCacheType,
		mock.AnythingOfType("Product"),
		time.Hour).
		Return(fmt.Errorf("error setting cache"))
	mProductRepo.On("SetStatus", mock.AnythingOfType("int"),
		domain.ExpiredProduct).Return(nil)
	_, err := interactor.GetUserAds(domain.Ad{UserID: 123})

	assert.NoError(t, err)
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserAdsErrorGetAds(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour, 2)
	productParams := domain.ProductParams{Categories: []int{2020}, Limit: 2}

	testTime := time.Now().Add(time.Hour * 24)
	product := domain.Product{Config: productParams,
		ExpiredAt: testTime, Status: domain.ActiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mLogger.On("LogErrorGettingUserAdsData", mock.Anything, mock.Anything)

	mAdRepo.On("GetUserAds", mock.AnythingOfType("int"),
		mock.AnythingOfType("ProductParams")).Return(domain.Ads{}, fmt.Errorf("err"))
	_, err := interactor.GetUserAds(domain.Ad{UserID: 123})

	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}

func TestGetUserAdsNotEnoughAds(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour, 2)
	productParams := domain.ProductParams{Categories: []int{2020}, Limit: 2}

	testTime := time.Now().Add(time.Hour * 24)
	product := domain.Product{Config: productParams,
		ExpiredAt: testTime, Status: domain.ActiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mLogger.On("LogNotEnoughAds", mock.Anything)

	mAdRepo.On("GetUserAds", mock.AnythingOfType("int"),
		mock.AnythingOfType("ProductParams")).Return(domain.Ads{domain.Ad{}}, nil)
	_, err := interactor.GetUserAds(domain.Ad{UserID: 123})

	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
