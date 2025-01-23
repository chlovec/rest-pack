package db

import (
	"database/sql"
	"log"
)

func NewSqlStorage(dbConnString string, dbDriver string, logger *log.Logger) (*sql.DB, error) {
	db, err := sql.Open(dbDriver, dbConnString)
	if err != nil {
		logger.Fatal(err)
	}

	return db, err
}

func InitDB(dbConnString string, dbDriver string, logger *log.Logger) (*sql.DB, error) {
	db, err := NewSqlStorage(dbConnString, dbDriver, logger)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}