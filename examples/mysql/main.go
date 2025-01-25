package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chlovec/rest-pack/db/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("examples/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Build the connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)


	mysqlDB, err := db.NewSqlStorage("mysql", dsn)
	if err != nil {
		log.Fatalf("error instantiating db: \v%v", err)
	}
	err = db.InitDB(mysqlDB)
	if err != nil {
		log.Fatalf("error connecting to db %v", err)
	}
	runQuery(mysqlDB)
}

func runQuery(db *sql.DB) {
	query := "SELECT id, email, lastName, firstName FROM users"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("error running query on db: %v", err)
	}
	defer rows.Close()

	// Get column names from the query result
	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("error fetching column names: %v", err)
	}

	// Print the column names
	fmt.Printf("Columns: %v\n", columns)

	// Iterate over the rows
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		// Create a slice of pointers to interface{} for Scan
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatalf("error scanning row: %v", err)
		}

		// Create a map of column names to values for better readability
		rowMap := make(map[string]interface{})
		for i, colName := range columns {
			val := values[i]

			// Handle NULL values gracefully
			var v interface{}
			switch b := val.(type) {
			case []byte:
				v = string(b) // Convert []byte to string
			default:
				v = b
			}
			rowMap[colName] = v
		}

		// Print the row as a map
		fmt.Printf("Row: %v\n", rowMap)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Fatalf("error iterating through rows: %v", err)
	}
}