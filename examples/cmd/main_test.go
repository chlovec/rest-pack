package main

import (
	"database/sql"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInitServer(t *testing.T) {
	addr := ":8080"
	logger := log.Default()

	// Create a mock database connection
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer mockDB.Close()

	// Mock the PingContext behavior to succeed
	mock.ExpectPing()

	t.Run("Success", func(t *testing.T) {
		// Override sqlOpen to return the mock database
		mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
			return mockDB, nil
		}

		apiServer, err := initServer(mockSQLOpen, addr, "testDsn", 0, logger)

		assert.NoError(t, err)
		assert.NotNil(t, apiServer)
	})

	t.Run("DB Error", func(t *testing.T) {
		// Override sqlOpen to return the mock database
		dbError := "DB Error"
		mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
			return nil, errors.New(dbError)
		}

		apiServer, err := initServer(mockSQLOpen, addr, "testDsn", 0, logger)

		assert.NotNil(t, err)
		assert.Equal(t, dbError, err.Error())
		assert.Nil(t, apiServer)
	})
}