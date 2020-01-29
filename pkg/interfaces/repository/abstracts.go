package repository

import (
	"io"
	"time"
)

// HTTPRequest interface represents the request that is going to be sent via HTTP
type HTTPRequest interface {
	GetMethod() string
	SetMethod(string) HTTPRequest
	GetPath() string
	SetPath(string) HTTPRequest
	GetBody() interface{}
	SetBody(interface{}) HTTPRequest
	GetHeaders() map[string][]string
	SetHeaders(map[string]string) HTTPRequest
	GetQueryParams() map[string][]string
	SetQueryParams(map[string]string) HTTPRequest
	GetTimeOut() time.Duration
	SetTimeOut(int) HTTPRequest
}

// HTTPHandler implements HTTP handler operations
type HTTPHandler interface {
	Send(HTTPRequest) (interface{}, error)
	NewRequest() HTTPRequest
}

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
