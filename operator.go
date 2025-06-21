package mysqlPool

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *Pool) Query(query string, params ...interface{}) (*sql.Rows, error) {
	if db.db == nil {
		return nil, db.Logger.Error(nil, "Database connection is not available")
	}

	startTime := time.Now()
	rows, err := db.db.Query(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		db.Logger.Info(fmt.Sprintf("Slow Query %s", duration))
	}

	return rows, err
}

func (db *Pool) Exec(query string, params ...interface{}) (sql.Result, error) {
	if db.db == nil {
		return nil, db.Logger.Error(nil, "Database connection is not available")
	}

	startTime := time.Now()
	result, err := db.db.Exec(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		db.Logger.Info(fmt.Sprintf("Slow Query %s", duration))
	}

	return result, err
}

// * private method
func (q *builder) query(query string, params ...interface{}) (*sql.Rows, error) {
	if q.db == nil {
		q.logger.Error(nil, "Database connection is not available")
		return nil, fmt.Errorf("database connection is not available")
	}

	startTime := time.Now()
	rows, err := q.db.Query(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		q.logger.Info(fmt.Sprintf("Slow Query %s", duration))
	}

	return rows, err
}

// * private method
func (q *builder) exec(query string, params ...interface{}) (sql.Result, error) {
	if q.db == nil {
		return nil, q.logger.Error(nil, "Database connection is not available")
	}

	startTime := time.Now()
	result, err := q.db.Exec(query, params...)
	duration := time.Since(startTime)

	if duration > 20*time.Millisecond {
		q.logger.Info(fmt.Sprintf("Slow Query %s", duration))
	}

	return result, err
}
