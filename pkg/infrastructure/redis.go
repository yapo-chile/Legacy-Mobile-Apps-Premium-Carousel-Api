package infrastructure

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/repository"
)

// RedisHandler handler for the request made to redis
type RedisHandler struct {
	Client *redis.Client
	Logger loggers.Logger
}

// NewRedisHandler constructor for RedisHandler
func NewRedisHandler(address, password string, db int, logger loggers.Logger) *RedisHandler {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password, // no password set
		DB:       db,       // use default DB
	})
	return &RedisHandler{
		Client: client,
		Logger: logger,
	}
}

// HGet gets the result of a HGET command with the given key/field
func (r *RedisHandler) HGet(key, field string) (string, bool) {
	cmdResult := r.Client.HGet(key, field)
	if err := cmdResult.Err(); err != nil {
		r.Logger.Error("redisError: %+v\n", err)
		return "", false
	}
	if result, err := cmdResult.Result(); err == nil {
		return result, true
	}
	r.Logger.Error("redisError: key(%s) field(%s)\n", key, field)
	return "", false
}

// HGetAll gets all result of a HGETALL command with the given key
func (r *RedisHandler) HGetAll(key string) (map[string]string, bool) {
	cmdResult := r.Client.HGetAll(key)
	if err := cmdResult.Err(); err != nil {
		r.Logger.Error("redisError: %+v\n", err)
		return map[string]string{}, false
	}
	if result, err := cmdResult.Result(); err == nil {
		r.Logger.Debug("redisResult: %+v\n", result)
		return result, true
	}
	r.Logger.Error("redisError: key(%s)\n", key)
	return map[string]string{}, false
}

// Get gets the result of a GET command with the given key
func (r *RedisHandler) Get(key string) (repository.RedisResult, error) {
	result := r.Client.Get(key)
	err := result.Err()
	if err == redis.Nil {
		return result, fmt.Errorf("KEY_NOT_FOUND: %s", key)
	}
	return result, err
}

// Set sets a value in redis with the given key
func (r *RedisHandler) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(key, value, expiration).Err()
}

// Del deletes the given key from the database in redis
func (r *RedisHandler) Del(key string) error {
	return r.Client.Del(key).Err()
}
