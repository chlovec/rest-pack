package db

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Test for successful initialization
func TestInitDB_Success(t *testing.T) {
	// Create a mock database connection
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Mock the PingContext behavior to succeed
	mock.ExpectPing()

	// Override sqlOpen to return the mock database
	mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}

	// Call the function under test with the mocked sqlOpen
	db, err := InitDB(mockSQLOpen, "mockDataSource")
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// Test for sql.Open error (failure to open the database)
func TestInitDB_OpenError(t *testing.T) {
	// Override sqlOpen to simulate an error
	mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, errors.New("failed to open database")
	}

	// Call the function under test
	db, err := InitDB(mockSQLOpen, "mockDataSource")

	// Assert that an error is returned and db is nil
	assert.Error(t, err)
	assert.Nil(t, db)
}

// Test for PingContext error (failure during ping)
func TestInitDB_PingError(t *testing.T) {
	// Create a mock database connection
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Mock the PingContext behavior to fail
	mock.ExpectPing().WillReturnError(errors.New("ping failed"))

	// Override sqlOpen to return the mock database
	mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
		return mockDB, nil
	}

	// Call the function under test
	db, err := InitDB(mockSQLOpen, "mockDataSource")

	// Assert that error is returned and db is nil
	assert.Error(t, err)
	assert.Nil(t, db)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
