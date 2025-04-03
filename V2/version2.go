package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// CreateDBConnection opens a connection to the PostgreSQL database
func CreateDBConnection() {
	var err error
	// Database connection string (replace with your credentials)
	connStr := "user=youruser password=yourpassword dbname=inventory_system sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	fmt.Println("Database connection established!")
}

// CloseDBConnection closes the database connection
func CloseDBConnection() {
	if db != nil {
		err := db.Close()
		if err != nil {
			log.Fatal("Error closing the database connection:", err)
		}
		fmt.Println("Database connection closed.")
	}
}

// CreateProductsTable creates the Products table if it doesn't exist
func CreateProductsTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Products (
			product_id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			price DECIMAL
		);
	`)
	if err != nil {
		log.Fatal("Error creating Products table:", err)
	}
	fmt.Println("Products table created (or already exists).")
}

// CreateInventoryTable creates the Inventory table if it doesn't exist
func CreateInventoryTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS Inventory (
			inventory_id SERIAL PRIMARY KEY,
			product_id INT REFERENCES Products(product_id),
			quantity INT
		);
	`)
	if err != nil {
		log.Fatal("Error creating Inventory table:", err)
	}
	fmt.Println("Inventory table created (or already exists).")
}

// CreateStockMovementsTable creates the StockMovements table if it doesn't exist
func CreateStockMovementsTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS StockMovements (
			movement_id SERIAL PRIMARY KEY,
			product_id INT REFERENCES Products(product_id),
			movement_type VARCHAR(50), -- 'stock-in', 'sale', 'manual-removal'
			quantity INT,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Error creating StockMovements table:", err)
	}
	fmt.Println("StockMovements table created (or already exists).")
}

func main() {
	// Create DB connection
	CreateDBConnection()
	defer CloseDBConnection()

	// Create tables if they don't exist
	CreateProductsTable()
	CreateInventoryTable()
	CreateStockMovementsTable()
}
