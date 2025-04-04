package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

// Db will be used globally to hold the database connection
var Db *sql.DB

// CreateDBConnection opens a connection to the PostgreSQL database
func CreateDBConnection(connStr string) (*sql.DB, error) {
	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Verify the connection
	err = Db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	fmt.Println("PostgreSQL database connection established!")
	return Db, nil
}

// CloseDBConnection closes the database connection
func CloseDBConnection() {
	if Db != nil {
		err := Db.Close()
		if err != nil {
			log.Fatal("Error closing the database connection:", err)
		}
		fmt.Println("PostgreSQL database connection closed.")
	}
}
