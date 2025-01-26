package db

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Test for successful initialization
func TestInitDB(t *testing.T) {
	// Define test cases
	testCases := []struct {
		driverName string
		timeout    int
	}{
		{"mysql", 0},    // MySQL with no timeout
		{"postgres", 5}, // Postgres with 5 seconds timeout
		{"pgx", 10},     // PGX with 10 seconds timeout
		{"sqlite3", 2},  // SQLite with 2 seconds timeout
	}
	for _, tc := range testCases {
		timeout := time.Duration(tc.timeout) * time.Second

		t.Run(fmt.Sprintf("success_with_%s_and_timeout_of_%d_secs", tc.driverName, tc.timeout), func(t *testing.T) {
			// Create a mock database connection
			mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			assert.NoError(t, err)
			defer mockDB.Close()

			// Mock the PingContext behavior to succeed
			mock.ExpectPing()

			// Override sqlOpen to return the mock database
			mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
				return mockDB, nil
			}

			// Call the function under test with the mocked sqlOpen
			db, err := InitDB(mockSQLOpen, tc.driverName, "mockDataSource", timeout)
			assert.NoError(t, err)
			assert.NotNil(t, db)

			// Ensure all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})

		t.Run(fmt.Sprintf("open_error_with_%s_and_timeout_of_%d_secs", tc.driverName, tc.timeout), func(t *testing.T) {
			// Override sqlOpen to simulate an error
			mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
				return nil, errors.New("failed to open database")
			}

			// Call the function under test
			db, err := InitDB(mockSQLOpen, tc.driverName, "mockDataSource", timeout)

			// Assert that an error is returned and db is nil
			assert.Error(t, err)
			assert.Nil(t, db)
		})

		t.Run(fmt.Sprintf("ping_error_with_%s_and_timeout_of_%d_secs", tc.driverName, tc.timeout), func(t *testing.T) {
			// Create a mock database connection
			mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			assert.NoError(t, err)
			defer mockDB.Close()

			// Mock the PingContext behavior to fail
			mock.ExpectPing().WillReturnError(errors.New("ping failed"))

			// Override sqlOpen to return the mock database
			mockSQLOpen := func(driverName, dataSourceName string) (*sql.DB, error) {
				return mockDB, nil
			}

			// Call the function under test
			db, err := InitDB(mockSQLOpen, "mysql", "mockDataSource", 10*time.Second)

			// Assert that error is returned and db is not nil
			assert.Error(t, err)
			assert.NotNil(t, db)

			// Ensure all expectations were met
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
