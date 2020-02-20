package repository

import (
	"fmt"
	"testing"
	"time"

	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/domain"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRedisResult struct {
	mock.Mock
}

func (m *mockRedisResult) Bytes() ([]byte, error) { // nolint
	args := m.Called()
	return args.Get(0).([]byte), args.Error(1)
}

type mockRedis struct {
	mock.Mock
}

func (m *mockRedis) HGetAll(key string) (map[string]string, bool) {
	args := m.Called(key)
	return args.Get(0).(map[string]string), args.Bool(1)
}

func (m *mockRedis) HGet(key, field string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *mockRedis) Set(key string, value interface{}, expiration time.Duration) error {
	args := m.Called(key, value, expiration)
	return args.Error(0)
}

func (m *mockRedis) Get(key string) (RedisResult, error) {
	args := m.Called(key)
	return args.Get(0).(RedisResult), args.Error(1)
}

func (m *mockRedis) Del(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func TestNewCacheRepository(t *testing.T) {
	m := &mockRedis{}
	expected := cacheRepository{
		handler:           m,
		defaultExpiration: time.Hour,
	}
	result := NewCacheRepository(m, "", time.Hour)
	assert.Equal(t, &expected, result)
	m.AssertExpectations(t)
}

func TestGetCache(t *testing.T) {
	m := &mockRedis{}
	mResult := &mockRedisResult{}
	repo := cacheRepository{
		handler:           m,
		defaultExpiration: time.Hour,
	}
	mResult.On("Bytes").Return([]byte{}, nil)
	m.On("Get", mock.AnythingOfType("string")).Return(mResult, nil)
	result, err := repo.GetCache(`some-key`, usecases.ProductCacheType)
	assert.NoError(t, err)
	assert.Equal(t, []byte{}, result)
	m.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestGetCacheErrorInRedis(t *testing.T) {
	m := &mockRedis{}
	mResult := &mockRedisResult{}
	repo := cacheRepository{
		handler:           m,
		defaultExpiration: time.Hour,
	}
	m.On("Get", mock.AnythingOfType("string")).Return(mResult, fmt.Errorf("err"))
	_, err := repo.GetCache(`some-key`, usecases.ProductCacheType)
	assert.Error(t, err)
	m.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestSetCacheOK(t *testing.T) {
	m := &mockRedis{}
	mResult := &mockRedisResult{}
	repo := cacheRepository{
		handler:           m,
		defaultExpiration: time.Hour,
	}
	m.On("Set", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8"),
		mock.AnythingOfType("time.Duration")).Return(nil)
	err := repo.SetCache(`some-key`, usecases.MinifiedAdDataType, domain.Ad{}, 0)
	assert.NoError(t, err)
	m.AssertExpectations(t)
	mResult.AssertExpectations(t)
}

func TestMinifyCacheOK(t *testing.T) {
	m := &mockRedis{}
	mResult := &mockRedisResult{}
	repo := cacheRepository{
		handler:           m,
		defaultExpiration: time.Hour,
	}
	ad := domain.Ad{
		ID:         "123",
		UserID:     "123",
		CategoryID: "2020",
		Price:      1111,
		Currency:   "USD",
	}
	expected := map[string]interface{}{
		"ID":         ad.ID,
		"UserID":     ad.UserID,
		"CategoryID": ad.CategoryID,
		"Price":      ad.Price,
		"Currency":   ad.Currency,
	}
	res := repo.minifyCache(usecases.MinifiedAdDataType, ad)
	assert.Equal(t, expected, res)
	res2 := repo.minifyCache(usecases.ProductCacheType, "test2")
	assert.Equal(t, "test2", res2)
	m.AssertExpectations(t)
	mResult.AssertExpectations(t)
}
