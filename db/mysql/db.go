package db

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func InitDB(sqlOpen func(driverName, dataSourceName string) (*sql.DB, error),dataSourceName string) (*sql.DB, error) {
	db, err := sqlOpen("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(context.Background())
	return db, err
}
