package usecases

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
)

type mockProductRepo struct {
	mock.Mock
}

func (m *mockProductRepo) GetUserActiveProduct(userID string,
	productType ProductType) (Product, error) {
	args := m.Called(userID, productType)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockProductRepo) GetUserProducts(email string, page int) ([]Product, int, int, error) {
	args := m.Called(email, page)
	return args.Get(0).([]Product), args.Int(1), args.Int(2), args.Error(3)
}

func (m *mockProductRepo) AddUserProduct(userID, email, comment string, productType ProductType,
	expiredAt time.Time, config CpConfig) (Product, error) {
	args := m.Called(userID, email, comment, productType, expiredAt, config)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockProductRepo) GetUserProductsTotal(email string) (total int) {
	args := m.Called(email)
	return args.Int(0)
}

func (m *mockProductRepo) GetUserProductByID(userProductID int) (Product, error) {
	args := m.Called(userProductID)
	return args.Get(0).(Product), args.Error(1)
}

func (m *mockProductRepo) SetConfig(userProductID int, config CpConfig) error {
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
	status ProductStatus) error {
	args := m.Called(userProductID, status)
	return args.Error(0)
}

type mockAdRepo struct {
	mock.Mock
}

func (m *mockAdRepo) GetUserAds(userID string,
	cpConfig CpConfig) (domain.Ads, error) {
	args := m.Called(userID, cpConfig)
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

func (m *mockgetUserAdsLogger) LogWarnGettingCache(userID string, err error) {
	m.Called(userID, err)
}

func (m *mockgetUserAdsLogger) LogWarnSettingCache(userID string, err error) {
	m.Called(userID, err)
}

func (m *mockgetUserAdsLogger) LogInfoActiveProductNotFound(userID string, product Product) {
	m.Called(userID, product)
}

func (m *mockgetUserAdsLogger) LogInfoProductExpired(userID string, product Product) {
	m.Called(userID, product)
}

func (m *mockgetUserAdsLogger) LogErrorGettingUserAdsData(userID string, err error) {
	m.Called(userID, err)
}

func TestGetUserAdsOkWithoutCache(t *testing.T) {
	mProductRepo := &mockProductRepo{}
	mAdRepo := &mockAdRepo{}
	mCacheRepo := &mockCacheRepo{}
	mLogger := &mockgetUserAdsLogger{}
	interactor := MakeGetUserAdsInteractor(mAdRepo, mProductRepo,
		mCacheRepo, mLogger, time.Hour)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}
	tAds := domain.Ads{
		{ID: "1", Subject: "Mi auto", UserID: "123"},
		{ID: "2", Subject: "Mi auto 2", UserID: "123"},
	}
	testTime := time.Now().Add(time.Hour * 24)
	product := Product{Config: cpConfig, UserID: "123",
		ExpiredAt: testTime, Status: ActiveProduct}
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return([]byte{}, fmt.Errorf("cache not found"))
	mLogger.On("LogWarnGettingCache", mock.Anything, mock.Anything)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)
	mProductRepo.On("GetUserActiveProduct", mock.AnythingOfType("string"),
		PremiumCarousel).Return(product, nil)
	mAdRepo.On("GetUserAds", mock.AnythingOfType("string"),
		mock.AnythingOfType("CpConfig")).Return(tAds, nil)
	mCacheRepo.On("SetCache", "user:123:PREMIUM_CAROUSEL",
		ProductCacheType,
		product,
		time.Hour).
		Return(fmt.Errorf("error setting cache"))
	ads, err := interactor.GetUserAds(domain.Ad{UserID: "123"})
	expected := tAds
	assert.NoError(t, err)
	assert.Equal(t, expected, ads)
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
		mCacheRepo, mLogger, time.Hour)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}
	tAds := domain.Ads{
		{ID: "1", Subject: "Mi auto", UserID: "123"},
		{ID: "2", Subject: "Mi auto 2", UserID: "123"},
	}
	testTime := time.Now().Add(time.Hour * 24)
	product := Product{Config: cpConfig,
		ExpiredAt: testTime, Status: ActiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)

	mAdRepo.On("GetUserAds", mock.AnythingOfType("string"),
		mock.AnythingOfType("CpConfig")).Return(tAds, nil)
	ads, err := interactor.GetUserAds(domain.Ad{UserID: "123"})
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
		mCacheRepo, mLogger, time.Hour)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}
	testTime := time.Now().Add(time.Hour * 24)
	product := Product{Config: cpConfig,
		ExpiredAt: testTime, Status: InactiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mLogger.On("LogInfoActiveProductNotFound", mock.Anything, mock.Anything)

	interactor.GetUserAds(domain.Ad{UserID: "123"})
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
		mCacheRepo, mLogger, time.Hour)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}
	testTime := time.Now().Add(time.Hour * -24)
	product := Product{Config: cpConfig,
		ExpiredAt: testTime, UserID: "123", Status: ActiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mLogger.On("LogInfoProductExpired", mock.Anything, mock.Anything)
	mLogger.On("LogWarnSettingCache", mock.Anything, mock.Anything)

	product.Status = ExpiredProduct
	mCacheRepo.On("SetCache", "user:123:PREMIUM_CAROUSEL",
		ProductCacheType,
		mock.AnythingOfType("Product"),
		time.Hour).
		Return(fmt.Errorf("error setting cache"))
	mProductRepo.On("SetStatus", mock.AnythingOfType("int"),
		ExpiredProduct).Return(nil)
	_, err := interactor.GetUserAds(domain.Ad{UserID: "123"})

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
		mCacheRepo, mLogger, time.Hour)
	cpConfig := CpConfig{Categories: []int{2020}, Limit: 2}

	testTime := time.Now().Add(time.Hour * 24)
	product := Product{Config: cpConfig,
		ExpiredAt: testTime, Status: ActiveProduct}
	productBytes, _ := json.Marshal(product)
	mCacheRepo.On("GetCache", mock.AnythingOfType("string"), ProductCacheType).
		Return(productBytes, nil)
	mLogger.On("LogErrorGettingUserAdsData", mock.Anything, mock.Anything)

	mAdRepo.On("GetUserAds", mock.AnythingOfType("string"),
		mock.AnythingOfType("CpConfig")).Return(domain.Ads{}, fmt.Errorf("err"))
	_, err := interactor.GetUserAds(domain.Ad{UserID: "123"})

	assert.Error(t, err)
	mProductRepo.AssertExpectations(t)
	mAdRepo.AssertExpectations(t)
	mCacheRepo.AssertExpectations(t)
	mLogger.AssertExpectations(t)
}
