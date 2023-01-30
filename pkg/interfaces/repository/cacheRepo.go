package repository

import (
	"encoding/json"
	"strings"
	"time"
	"log"

	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/domain"
	"gitlab.com/yapo_team/legacy/mobile-apps/premium-carousel-api/pkg/usecases"
)

// cacheRepository allows get cached request responses using redis handler
type cacheRepository struct {
	handler           Redis
	prefix            string
	defaultExpiration time.Duration
}

// NewCacheRepository returns an instance of cacheRepository
func NewCacheRepository(handler Redis, prefix string, defaultExpiration time.Duration) usecases.CacheRepository {
	return &cacheRepository{
		handler:           handler,
		prefix:            prefix,
		defaultExpiration: defaultExpiration,
	}
}

// makeRedisKey generates key for redis
func (repo *cacheRepository) makeRedisKey(key string, cacheType usecases.CacheType) string {
	return strings.Join([]string{key, string(cacheType)}, ":")
}

// GetCache returns the response of a cached request
func (repo *cacheRepository) GetCache(key string, cacheType usecases.CacheType) ([]byte, error) {
	k := repo.makeRedisKey(key, cacheType)
	res, err := repo.handler.Get(k)
	if err != nil {
		return nil, err
	}
	return res.Bytes()
}

// SetCache saves the response of request in redis
func (repo *cacheRepository) SetCache(key string, cacheType usecases.CacheType,
	data interface{}, expiration time.Duration) error {
	log.Printf("Set cache")
	if expiration <= 0 {
		expiration = repo.defaultExpiration
	}
	log.Printf("After set expiration")
	k := repo.makeRedisKey(key, cacheType)
	log.Printf("After make redis key")
	data = repo.minifyCache(cacheType, data)
	log.Printf("After minify redis key")
	bytes, _ := json.Marshal(data) // nolint
	log.Printf("After marshal data")
	return repo.handler.Set(k, bytes, expiration)
}

// minifyCache tries to reduce known cache types
func (repo *cacheRepository) minifyCache(cacheType usecases.CacheType,
	data interface{}) interface{} {
	switch cacheType {
	case usecases.MinifiedAdDataType:
		ad := data.(domain.Ad)
		return map[string]interface{}{
			"ID":         ad.ID,
			"UserID":     ad.UserID,
			"CategoryID": ad.CategoryID,
			"Price":      ad.Price,
			"Currency":   ad.Currency,
		}
	default:
		return data
	}
}
