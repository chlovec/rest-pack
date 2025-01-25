package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestNewSqlStorage(t *testing.T) {
	// Create a mock database connection
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Mock driver name and data source name
	driverName := "mysql"
	dataSourceName := "user:password@tcp(127.0.0.1:3306)/mockdb"

	// Call the function being tested
	dbConn, err := NewSqlStorage(driverName, dataSourceName)

	// Assertions
	assert.NoError(t, err)        // Ensure no error occurred
	assert.NotNil(t, dbConn)      // Ensure the returned connection is not nil

	// Close the database connection
	err = dbConn.Close()
	assert.NoError(t, err) // Ensure no error occurred when closing
}

func TestInitDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing()

	err = InitDB(db)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}
