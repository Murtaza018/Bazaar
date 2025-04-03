package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var Db *sql.DB

// CreateDBConnection opens a connection to the SQLite database
func CreateDBConnection() {
	var err error
	// Database connection string (replace with your file path)
	connStr := "Bazaar_v1.db"
	Db, err = sql.Open("sqlite3", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	fmt.Println("Database connection established!")
}

// CloseDBConnection closes the database connection
func CloseDBConnection() {
	if Db != nil {
		err := Db.Close()
		if err != nil {
			log.Fatal("Error closing the database connection:", err)
		}
		fmt.Println("Database connection closed.")
	}
}
