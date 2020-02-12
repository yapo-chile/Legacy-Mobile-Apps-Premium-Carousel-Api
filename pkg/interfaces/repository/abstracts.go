package repository

import (
	"encoding/json"
	"io"
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

// Config contains all info of configured
type Config interface {
	Get(string) string
}

type SearchResult interface {
	GetResults() (results []json.RawMessage)
	TotalHits() int64
}

type Query interface {
	// Source returns the JSON-serializable query request.
	Source() (interface{}, error)
}

type Search interface {
	NewMultiMatchQuery(text interface{}, typ string, fields ...string) Query
	NewTermQuery(name string, value interface{}) Query
	NewRangeQuery(name string, from, to int) Query
	NewFunctionScoreQuery(query Query, boost float64, boostMode string, random bool) Query
	NewBoolQuery(must []Query, mustNot []Query) Query
	NewIDsQuery(ids ...string) Query
	NewCategoryFilter(categoryIDs ...int) Query
	GetDoc(index string, id string) (json.RawMessage, error)
	Search(index string, query Query, from, size int) (SearchResult, error)
}
