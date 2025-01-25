package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func NewSqlStorage(driverName string, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return db, err
}

func InitDB(db *sql.DB) error {
	return db.Ping()
}