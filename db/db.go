package db

import (
	"context"
	"database/sql"
	"time"
)

func InitDB(sqlOpen func(driverName, dataSourceName string) (*sql.DB, error),driverName string, dataSourceName string, timeout time.Duration) (*sql.DB, error) {
	if timeout == 0 {
		timeout = 2*time.Second
	}

	db, err := sqlOpen(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	err = db.PingContext(ctx)
	return db, err
}