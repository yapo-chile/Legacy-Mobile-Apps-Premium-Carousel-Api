package repository

import (
	"encoding/json"
	"io"
	"time"
)

// DbHandler represents a database connection handler
// it provides basic database capabilities
// after its use, the connection with the database must be closed
type DbHandler interface {
	io.Closer
	Insert(statement string, params ...interface{}) error
	Update(statement string, params ...interface{}) error
	Query(statement string, params ...interface{}) (DbResult, error)
}

// DbResult represents a database query result rows
// after its use, the Close() method must be invoked
// to ensure that the database connection used to perform the query
// returns to the connection pool to be used again
type DbResult interface {
	io.Closer
	Scan(dest ...interface{})
	Next() bool
}

// Redis implements Redis functions
type Redis interface {
	HGetAll(key string) (map[string]string, bool)
	HGet(key, field string) (string, bool)
	Set(key string, values interface{}, expiration time.Duration) error
	Get(key string) (RedisResult, error)
	Del(key string) error
}

// RedisResult interface for a result obtained from executing a get command in redis
type RedisResult interface {
	Bytes() ([]byte, error)
}

// Config contains all info of configured
type Config interface {
	Get(string) string
}

// SearchResult interface for a result obtained from search repository
type SearchResult interface {
	GetResults() (results []json.RawMessage)
	TotalHits() int64
}

// Query interface for a query request using search repository
type Query interface {
	Source() (interface{}, error)
}

// Search allows search over ads documents using external repository
type Search interface {
	NewMultiMatchQuery(text interface{}, typ string, fields ...string) Query
	NewTermQuery(name string, value interface{}) Query
	NewRangeQuery(name string, from, to int) Query
	NewFunctionScoreQuery(query Query, boost float64, boostMode string, random bool) Query
	NewBoolQuery(must, mustNot, should []Query) Query
	NewIDsQuery(ids ...string) Query
	NewCategoryFilter(categoryIDs ...int) Query
	GetDoc(index string, id string) (json.RawMessage, error)
	Search(index string, query Query, from, size int) (SearchResult, error)
}

// KafkaProducer allows send messages to kafka
type KafkaProducer interface {
	SendMessage(topic string, message []byte) error
	io.Closer
}
