package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/loggers"
	"github.mpi-internal.com/Yapo/premium-carousel-api/pkg/interfaces/repository"
)

// PgsqlHandler allows connection with postgres database
type PgsqlHandler struct {
	Conn *sql.DB
}

// Healthcheck verifies a connection to the database is still alive
func (handler *PgsqlHandler) Healthcheck() bool {
	return handler.Conn.Ping() == nil
}

// Close closes db connection
func (handler *PgsqlHandler) Close() error {
	return handler.Conn.Close()
}

// Insert executes an insert query in db
func (handler *PgsqlHandler) Insert(statement string, params ...interface{}) error {
	_, err := handler.Conn.Exec(statement, params...)
	return err
}

// Update executes an update query in db
func (handler *PgsqlHandler) Update(statement string, params ...interface{}) error {
	_, err := handler.Conn.Exec(statement, params...)
	return err
}

// Query executes a query that returns rows, typically a SELECT.
func (handler *PgsqlHandler) Query(statement string, params ...interface{}) (repository.DbResult, error) {
	rows, err := handler.Conn.Query(statement, params...)
	if err != nil {
		fmt.Println(err)
		return new(PgsqlRow), err
	}
	return PgsqlRow{
		Rows: rows,
	}, nil
}

// PgsqlRow represents the result of a query
type PgsqlRow struct {
	Rows *sql.Rows
}

// Scan copies the columns in the current row into the values pointed at by dest
func (r PgsqlRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

// Next prepares the next result row for reading with the Scan method
func (r PgsqlRow) Next() bool {
	return r.Rows.Next()
}

// Close closes query result
func (r PgsqlRow) Close() error {
	return r.Rows.Close()
}

// MakePgsqlHandler creates a new instance of pgsql connector
func MakePgsqlHandler(conf DatabaseConf, logger loggers.Logger) (*PgsqlHandler, error) {
	poolDb, err := sql.Open("postgres",
		fmt.Sprintf("host=%s dbname=%s port=%d sslmode=%s user=%s password=%s",
			conf.Host, conf.Dbname, conf.Port, conf.Sslmode, conf.DbUser, conf.DbPasswd),
	)

	if err != nil || poolDb == nil {
		logger.Error("Error on pool DB definition %+v\n", err)
		return nil, err
	}

	for i := 0; i < conf.ConnRetries; i++ {
		if err := poolDb.Ping(); err != nil {
			logger.Info("Connection attempt number %d failed. Error: %+v", i+1, err)
			time.Sleep(1 * time.Second)
		} else {
			poolDb.SetMaxIdleConns(conf.MaxIdle)
			poolDb.SetMaxOpenConns(conf.MaxOpen)

			return &PgsqlHandler{
				Conn: poolDb,
			}, nil
		}
	}
	return nil, fmt.Errorf("Max connection attemtps reached")
}
